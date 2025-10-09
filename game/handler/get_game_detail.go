package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/GameLaunchPad/game_management_project/service"
	"gorm.io/gorm"
)

func GetGameDetail(ctx context.Context, req *game.GetGameDetailRequest) (*game.GetGameDetailResponse, error) {
	// parameter validation
	if req.GameID <= 0 {
		return &game.GetGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid GameID"},
		}, nil
	}

	// get game detail from DAO
	gameDdl, newestVersionDdl, onlineVersionDdl, err := GameDao.GetGameDetail(ctx, uint64(req.GameID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &game.GetGameDetailResponse{
				BaseResp: &common.BaseResp{Code: "10001", Msg: "Game not found"},
			}, nil
		}
		return &game.GetGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Failed to get game detail: " + err.Error()},
		}, nil
	}

	// transform to response format
	gameDetail, err := service.ConvertDdlToDetailGame(gameDdl, newestVersionDdl, onlineVersionDdl)
	if err != nil {
		return &game.GetGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Failed to convert game data: " + err.Error()},
		}, nil
	}

	// construct response
	resp := &game.GetGameDetailResponse{
		GameDetail: gameDetail,
		BaseResp:   &common.BaseResp{Code: "200", Msg: "Success"},
	}

	return resp, nil
}
