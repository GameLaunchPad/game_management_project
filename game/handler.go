package main

import (
	"context"
	game "github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

// GameServiceImpl implements the last service interface defined in the IDL.
type GameServiceImpl struct{}

// GetGameList implements the GameServiceImpl interface.
func (s *GameServiceImpl) GetGameList(ctx context.Context, req *game.GetGameListRequest) (resp *game.GetGameListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetGameDetail implements the GameServiceImpl interface.
func (s *GameServiceImpl) GetGameDetail(ctx context.Context, req *game.GetGameDetailRequest) (resp *game.GetGameDetailResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateGameDraft implements the GameServiceImpl interface.
func (s *GameServiceImpl) UpdateGameDraft(ctx context.Context, req *game.UpdateGameDraftRequest) (resp *game.UpdateGameDraftResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateGameDetail implements the GameServiceImpl interface.
func (s *GameServiceImpl) CreateGameDetail(ctx context.Context, req *game.CreateGameDetailRequest) (resp *game.CreateGameDetailResponse, err error) {
	// TODO: Your code here...
	return
}

// ReviewGameVersion implements the GameServiceImpl interface.
func (s *GameServiceImpl) ReviewGameVersion(ctx context.Context, req *game.ReviewGameVersionRequest) (resp *game.ReviewGameVersionResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteGameDraft implements the GameServiceImpl interface.
func (s *GameServiceImpl) DeleteGameDraft(ctx context.Context, req *game.DeleteGameDraftRequest) (resp *game.DeleteGameDraftResponse, err error) {
	// TODO: Your code here...
	return
}
