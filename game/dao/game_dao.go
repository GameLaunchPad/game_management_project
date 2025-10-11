package dao

import (
	"context"
	"errors"
	"time"

	"github.com/GameLaunchPad/game_management_project/dal"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"gorm.io/gorm"
)

// GameWithVersionStatus is a struct to hold the result of a JOIN query
// between gp_game and gp_game_version.
type GameWithVersionStatus struct {
	ddl.GpGame
	Status int `gorm:"column:status"`
}

// CreateGame creates a new game and its initial version in a transaction.
func (d *gameDAO) CreateGame(ctx context.Context, game *ddl.GpGame, version *ddl.GpGameVersion) error {
	return dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. create gp_game record
		if err := tx.Create(game).Error; err != nil {
			return err
		}
		// 2. create gp_game_version record
		if err := tx.Create(version).Error; err != nil {
			return err
		}
		// 3. rollback newest_game_version_id
		if err := tx.Model(game).Update("newest_game_version_id", version.Id).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateGameVersionAndUpdateGame creates a new game version and updates the main game's newest_game_version_id.
func (d *gameDAO) CreateGameVersionAndUpdateGame(ctx context.Context, gameID uint64, version *ddl.GpGameVersion) error {
	return dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. check if the game exists
		var game ddl.GpGame
		if err := tx.First(&game, gameID).Error; err != nil {
			return err
		}

		// 2. create new gp_game_version record
		if err := tx.Create(version).Error; err != nil {
			return err
		}

		// 3. update newest_game_version_id in gp_game
		if err := tx.Model(&game).Update("newest_game_version_id", version.Id).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetGameList retrieves a paginated list of games with the status of their newest version.
func (d *gameDAO) GetGameList(ctx context.Context, filterText *string, pageNum, pageSize int) ([]*GameWithVersionStatus, int64, error) {
	var results []*GameWithVersionStatus
	var total int64

	// Start building the query on the gp_game table, aliased as 'g'
	db := dal.DB.WithContext(ctx).Model(&ddl.GpGame{}).Table("gp_game AS g")

	// Apply filter if provided
	if filterText != nil && *filterText != "" {
		db = db.Where("g.game_name LIKE ?", "%"+*filterText+"%")
	}

	// First, count the total number of records that match the filter
	// We need to perform the JOIN even for counting if the filter applies to the joined table in the future.
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (pageNum - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Now, perform the JOIN query to get the full data for the current page
	// JOIN gp_game_version (aliased as 'gv') on the newest_game_version_id
	// SELECT g.* (all columns from gp_game) and gv.status
	err := db.Select("g.*, gv.status").
		Joins("LEFT JOIN gp_game_version AS gv ON g.newest_game_version_id = gv.id").
		Order("g.modify_ts DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetGameDetail retrieves the main game info and its associated newest and online versions.
func (d *gameDAO) GetGameDetail(ctx context.Context, gameID uint64) (*ddl.GpGame, *ddl.GpGameVersion, *ddl.GpGameVersion, error) {
	var game ddl.GpGame
	var newestVersion *ddl.GpGameVersion
	var onlineVersion *ddl.GpGameVersion

	// 1. get the main game info
	if err := dal.DB.WithContext(ctx).First(&game, gameID).Error; err != nil {
		// if record not found, return nils
		return nil, nil, nil, err
	}

	// 2. get the newest game version
	if game.NewestGameVersionId != 0 {
		var nv ddl.GpGameVersion
		if err := dal.DB.WithContext(ctx).First(&nv, game.NewestGameVersionId).Error; err == nil {
			newestVersion = &nv
		}
	}

	// 3. get the online game version
	if game.OnlineGameVersionId != 0 {
		// if the online version is the same as the newest version, reuse it
		if newestVersion != nil && game.OnlineGameVersionId == newestVersion.Id {
			onlineVersion = newestVersion
		} else {
			var ov ddl.GpGameVersion
			if err := dal.DB.WithContext(ctx).First(&ov, game.OnlineGameVersionId).Error; err == nil {
				onlineVersion = &ov
			}
		}
	}

	return &game, newestVersion, onlineVersion, nil
}

// ReviewGameVersion updates a game version's status and potentially the main game's online version.
func (d *gameDAO) ReviewGameVersion(ctx context.Context, gameID, versionID uint64, newStatus int, reviewComment string) error {
	return dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. update gp_game_version status and review info
		updateData := map[string]interface{}{
			"status":         newStatus,
			"review_comment": reviewComment,
			"review_time":    time.Now().Unix(),
		}

		result := tx.Model(&ddl.GpGameVersion{}).Where("id = ? AND game_id = ?", versionID, gameID).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			// if no rows were affected, it means either the version_id or game_id is invalid
			return gorm.ErrRecordNotFound
		}

		// 2. if the new status is Published, update gp_game's online_game_version_id
		if newStatus == int(game.GameStatus_Published) {
			result = tx.Model(&ddl.GpGame{}).Where("id = ?", gameID).Update("online_game_version_id", versionID)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}

		return nil
	})
}

var ErrVersionIsNotDraft = errors.New("the newest version of the game is not a draft")

// DeleteGameDraft finds the newest version of a game, and if it's a draft, updates its status to Rejected.
func (d *gameDAO) DeleteGameDraft(ctx context.Context, gameID uint64) error {
	return dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. find the game record
		var gameRecord ddl.GpGame
		if err := tx.First(&gameRecord, gameID).Error; err != nil {
			return err
		}

		if gameRecord.NewestGameVersionId == 0 {
			// no versions exist for this game
			return gorm.ErrRecordNotFound
		}

		// 2. find the newest version record
		var newestVersion ddl.GpGameVersion
		if err := tx.First(&newestVersion, gameRecord.NewestGameVersionId).Error; err != nil {
			return err
		}

		// 3. check if the newest version is a draft
		if newestVersion.Status != int(game.GameStatus_Draft) {
			// if it's not a draft, cannot delete
			return ErrVersionIsNotDraft
		}

		// 4. delete (mark as Rejected) the draft version
		updateData := map[string]interface{}{
			"status": int(game.GameStatus_Rejected),
		}
		if err := tx.Model(&newestVersion).Updates(updateData).Error; err != nil {
			return err
		}

		return nil
	})
}
