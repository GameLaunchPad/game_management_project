package handler

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/dao"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/GameLaunchPad/game_management_project/service"
	"github.com/yitter/idgenerator-go/idgen"
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

	if req.GameDetail.GameID != 0 {
		return &game.CreateGameDetailResponse{
			BaseResp: &common.BaseResp{
				Code: "400",
				Msg:  "GameID must be 0 when creating a new game",
			},
		}, nil
	}

	// generate new IDs
	gameID := uint64(idgen.NextId())
	versionID := uint64(idgen.NextId())

	// Convert GameDetail and GameVersion to DDL structs
	gameVersionDdl, err := service.ConvertGameVersionToDdl(req.GameDetail.GameVersion)
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
