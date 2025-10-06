package dao

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dal"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
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
