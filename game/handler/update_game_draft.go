package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/GameLaunchPad/game_management_project/game/service"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

func UpdateGameDraft(ctx context.Context, req *game.UpdateGameDraftRequest) (*game.UpdateGameDraftResponse, error) {
	// param validation
	if req.GameDetail == nil || req.GameDetail.GameVersion == nil {
		return &game.UpdateGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid request: GameDetail or GameVersion is missing"},
		}, nil
	}
	// when updating a draft, GameID must be provided
	if req.GameDetail.GameID <= 0 {
		return &game.UpdateGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "GameID is required for updating a draft"},
		}, nil
	}

	// generate new version ID
	gameID := uint64(req.GameDetail.GameID)
	versionID := uint64(idgen.NextId())

	// convert GameVersion to DDL struct
	gameVersionDdl, err := service.ConvertGameVersionToDdl(req.GameDetail.GameVersion)
	if err != nil {
		return &game.UpdateGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid game version data: " + err.Error()},
		}, nil
	}
	gameVersionDdl.Id = versionID
	gameVersionDdl.GameId = gameID

	gameVersionDdl.Status = int(game.GameStatus_Draft)

	// call DAO to update the draft
	err = GameDao.UpdateGameDraft(ctx, gameID, gameVersionDdl)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &game.UpdateGameDraftResponse{
				BaseResp: &common.BaseResp{Code: "10001", Msg: "Game not found"},
			}, nil
		}
		return &game.UpdateGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Internal Server Error: " + err.Error()},
		}, nil
	}

	// construct success response
	return &game.UpdateGameDraftResponse{
		BaseResp: &common.BaseResp{Code: "200", Msg: "Success"},
	}, nil
}
