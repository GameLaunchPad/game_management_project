package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/dao"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

var GameDao dao.IGameDAO

func CreateGameDetail(ctx context.Context, req *game.CreateGameDetailRequest) (*game.CreateGameDetailResponse, error) {
	if req.GameDetail == nil || req.GameDetail.GameVersion == nil {
		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{
				Code: "400",
				Msg:  "Invalid request: GameDetail or GameVersion is missing",
			},
		}, nil
	}

	isUpdate := req.GameDetail.GameID != 0

	if isUpdate {
		return handleUpdateGame(ctx, req)
	} else {
		return handleCreateGame(ctx, req)
	}
}

func handleCreateGame(ctx context.Context, req *game.CreateGameDetailRequest) (*game.CreateGameDetailResponse, error) {
	reqVersion := req.GameDetail.GameVersion

	gameID := uint64(idgen.NextId())
	versionID := uint64(idgen.NextId())

	gameVersionDdl, err := ConvertGameVersionToDdl(reqVersion)
	if err != nil {
		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid game version data: " + err.Error()},
		}, nil
	}
	gameVersionDdl.Id = versionID
	gameVersionDdl.GameId = gameID

	gameDdl := &ddl.GpGame{
		Id:                  gameID,
		CpId:                uint64(req.GameDetail.CpID),
		GameName:            gameVersionDdl.GameName,
		GameIcon:            gameVersionDdl.GameIcon,
		HeaderImage:         gameVersionDdl.HeaderImage,
		NewestGameVersionId: versionID,
	}

	if err := GameDao.CreateGame(ctx, gameDdl, gameVersionDdl); err != nil {
		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Internal Server Error: " + err.Error()},
		}, nil
	}

	return &game.CreateGameDetailResponse{
		GameID:   int64(gameID),
		BaseResp: &common.BaseResp{Code: "200", Msg: "Success"},
	}, nil
}

func handleUpdateGame(ctx context.Context, req *game.CreateGameDetailRequest) (*game.CreateGameDetailResponse, error) {
	gameID := uint64(req.GameDetail.GameID)
	versionID := uint64(idgen.NextId())
	reqVersion := req.GameDetail.GameVersion

	gameVersionDdl, err := ConvertGameVersionToDdl(reqVersion)
	if err != nil {
		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid game version data: " + err.Error()},
		}, nil
	}
	gameVersionDdl.Id = versionID
	gameVersionDdl.GameId = gameID

	err = GameDao.CreateGameVersionAndUpdateGame(ctx, gameID, gameVersionDdl)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &game.CreateGameDetailResponse{
				BaseResp: &common.BaseResp{Code: "10001", Msg: "Game not found"},
			}, nil
		}

		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Internal Server Error: " + err.Error()},
		}, nil
	}

	return &game.CreateGameDetailResponse{
		GameID:   int64(gameID),
		BaseResp: &common.BaseResp{Code: "200", Msg: "Success"},
	}, nil
}
