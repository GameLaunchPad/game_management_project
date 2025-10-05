package dao

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dal"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
	"gorm.io/gorm"
)

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
