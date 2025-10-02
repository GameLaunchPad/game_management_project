package main

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/handler"
	cp_center "github.com/GameLaunchPad/game_management_project/kitex_gen/cp_center"
	game "github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

// GameServiceImpl implements the last service interface defined in the IDL.
type GameServiceImpl struct{}

// CpCenterServiceImpl implements the CpCenter service interface defined in the IDL.
type CpCenterServiceImpl struct{}

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

// CreateCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) CreateCPMaterial(ctx context.Context, req *cp_center.CreateCPMaterialRequest) (resp *cp_center.CreateCPMaterialResponse, err error) {
	// TODO: Your code here...
	return handler.NewCpCenterServiceImpl().CreateCPMaterial(ctx, req)
}

// UpdateCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) UpdateCPMaterial(ctx context.Context, req *cp_center.UpdateCPMaterialRequest) (resp *cp_center.UpdateCPMaterialResponse, err error) {
	// TODO: Your code here...
	return handler.NewCpCenterServiceImpl().UpdateCPMaterial(ctx, req)
}

// ReviewCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) ReviewCPMaterial(ctx context.Context, req *cp_center.ReviewCPMaterialRequest) (resp *cp_center.ReviewCPMaterialResponse, err error) {
	// TODO: Your code here...
	return handler.NewCpCenterServiceImpl().ReviewCPMaterial(ctx, req)
}

// GetCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) GetCPMaterial(ctx context.Context, req *cp_center.GetCPMaterialRequest) (resp *cp_center.GetCPMaterialResponse, err error) {
	// TODO: Your code here...
	return handler.NewCpCenterServiceImpl().GetCPMaterial(ctx, req)
}
