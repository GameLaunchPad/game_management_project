package dao

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dao/ddl"
)

// IGameDAO defines the interface for game data access operations.
type IGameDAO interface {
	CreateGame(ctx context.Context, game *ddl.GpGame, version *ddl.GpGameVersion) error
	UpdateGameDraft(ctx context.Context, gameID uint64, version *ddl.GpGameVersion) error
	GetGameList(ctx context.Context, filterText *string, pageNum, pageSize int) ([]*GameWithVersionStatus, int64, error)
	GetGameDetail(ctx context.Context, gameID uint64) (*ddl.GpGame, *ddl.GpGameVersion, *ddl.GpGameVersion, error)
	ReviewGameVersion(ctx context.Context, gameID, versionID uint64, newStatus int, reviewComment string) error
	DeleteGameDraft(ctx context.Context, gameID uint64) error
}
