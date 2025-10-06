package handler

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

// GetGameList handles the business logic for getting a list of games.
func GetGameList(ctx context.Context, req *game.GetGameListRequest) (*game.GetGameListResponse, error) {
	// parse parameters
	var filterText *string
	if req.IsSetFilter() && req.Filter.IsSetFilterText() {
		filterText = req.Filter.FilterText
	}

	pageNum := int(req.PageNum)
	if pageNum <= 0 {
		pageNum = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// get game list from DAO
	gamesDdl, total, err := GameDao.GetGameList(ctx, filterText, pageNum, pageSize)
	if err != nil {
		return &game.GetGameListResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Failed to get game list: " + err.Error()},
		}, nil
	}

	// transform to response format
	briefGames := make([]*game.BriefGame, 0, len(gamesDdl))
	for _, gameDdl := range gamesDdl {
		briefGame, err := ConvertDdlToBriefGame(gameDdl)
		if err != nil {
			return &game.GetGameListResponse{
				BaseResp: &common.BaseResp{Code: "500", Msg: "Failed to convert game data: " + err.Error()},
			}, nil
		}
		briefGames = append(briefGames, briefGame)
	}

	// construct response
	resp := &game.GetGameListResponse{
		GameList:   briefGames,
		TotalCount: int32(total),
		BaseResp:   &common.BaseResp{Code: "200", Msg: "Success"},
	}

	return resp, nil
}
