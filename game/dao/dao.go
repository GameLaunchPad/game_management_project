package dao

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dao/ddl"
)

// IGameDAO defines the interface for game data access operations.
type IGameDAO interface {
	CreateGame(ctx context.Context, game *ddl.GpGame, version *ddl.GpGameVersion) error
	CreateGameVersionAndUpdateGame(ctx context.Context, gameID uint64, version *ddl.GpGameVersion) error
}
