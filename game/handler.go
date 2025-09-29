package main

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/handler"
	game "github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

// GameServiceImpl implements the last service interface defined in the IDL.
type GameServiceImpl struct{}

// GetGameList implements the GameServiceImpl interface.
func (s *GameServiceImpl) GetGameList(ctx context.Context, req *game.GetGameListRequest) (resp *game.GetGameListResponse, err error) {
	return handler.GetGameList(ctx, req)
}

// GetGameDetail implements the GameServiceImpl interface.
func (s *GameServiceImpl) GetGameDetail(ctx context.Context, req *game.GetGameDetailRequest) (resp *game.GetGameDetailResponse, err error) {
	return handler.GetGameDetail(ctx, req)
}

// CreateGameDetail implements the GameServiceImpl interface.
func (s *GameServiceImpl) CreateGameDetail(ctx context.Context, req *game.CreateGameDetailRequest) (resp *game.CreateGameDetailResponse, err error) {
	return handler.CreateGameDetail(ctx, req)
}

// ReviewGameVersion implements the GameServiceImpl interface.
func (s *GameServiceImpl) ReviewGameVersion(ctx context.Context, req *game.ReviewGameVersionRequest) (resp *game.ReviewGameVersionResponse, err error) {
	return handler.ReviewGameVersion(ctx, req)
}

// DeleteGameDraft implements the GameServiceImpl interface.
func (s *GameServiceImpl) DeleteGameDraft(ctx context.Context, req *game.DeleteGameDraftRequest) (resp *game.DeleteGameDraftResponse, err error) {
	return handler.DeleteGameDraft(ctx, req)
}
