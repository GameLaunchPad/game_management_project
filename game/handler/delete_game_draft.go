package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/game/dao"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"gorm.io/gorm"
)

func DeleteGameDraft(ctx context.Context, req *game.DeleteGameDraftRequest) (*game.DeleteGameDraftResponse, error) {
	// --- 1. 参数校验 ---
	if req.GameID <= 0 {
		return &game.DeleteGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid GameID"},
		}, nil
	}

	// --- 2. 调用 DAO 层执行软删除 ---
	err := GameDao.DeleteGameDraft(ctx, uint64(req.GameID))
	if err != nil {
		// 检查是否是“记录未找到”的特定错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &game.DeleteGameDraftResponse{
				BaseResp: &common.BaseResp{Code: "10001", Msg: "Game not found"},
			}, nil
		}
		// 检查是否是“版本不是草稿”的特定错误
		if errors.Is(err, dao.ErrVersionIsNotDraft) {
			return &game.DeleteGameDraftResponse{
				BaseResp: &common.BaseResp{Code: "10003", Msg: err.Error()}, // 假设 10003 是业务错误码
			}, nil
		}
		// 其他所有数据库错误都归为内部错误
		return &game.DeleteGameDraftResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Internal Server Error: " + err.Error()},
		}, nil
	}

	// --- 3. 构建并返回成功的响应 ---
	resp := &game.DeleteGameDraftResponse{
		BaseResp: &common.BaseResp{Code: "200", Msg: "Success"},
	}
	return resp, nil
}
