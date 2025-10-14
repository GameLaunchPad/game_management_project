package dao

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dao/ddl"
)

type gameDAO struct{}

func (d *gameDAO) GetGameDetail(ctx context.Context, gameID uint64) (*ddl.GpGame, *ddl.GpGameVersion, *ddl.GpGameVersion, error) {
	// ToDo: 待实现
	return nil, nil, nil, nil
}

// NewGameDAO creates a new GameDAO.
func NewGameDAO() IGameDAO {
	return &gameDAO{}
}
