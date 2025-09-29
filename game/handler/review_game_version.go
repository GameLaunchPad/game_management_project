package handler

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

func ReviewGameVersion(ctx context.Context, req *game.ReviewGameVersionRequest) (*game.ReviewGameVersionResponse, error) {
	resp := &game.ReviewGameVersionResponse{
		BaseResp: &common.BaseResp{},
	}
	return resp, nil
}
