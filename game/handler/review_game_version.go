package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"gorm.io/gorm"
)

func ReviewGameVersion(ctx context.Context, req *game.ReviewGameVersionRequest) (*game.ReviewGameVersionResponse, error) {
	// --- 1. 参数校验 ---
	if req.GameID <= 0 || req.GameVersionID <= 0 {
		return &game.ReviewGameVersionResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid GameID or GameVersionID"},
		}, nil
	}

	// --- 2. 将审核结果 (RPC enum) 转换为数据库状态 (int) ---
	var newStatus int
	switch req.ReviewResult_ {
	case game.ReviewResult__Pass:
		newStatus = int(game.GameStatus_Published)
	case game.ReviewResult__Reject:
		newStatus = int(game.GameStatus_Rejected)
	default:
		return &game.ReviewGameVersionResponse{
			BaseResp: &common.BaseResp{Code: "400", Msg: "Invalid review result"},
		}, nil
	}

	// 注意：当前请求中没有 reviewComment 字段，我们暂时传入空字符串。
	// 这是一个未来可以优化的地方，可以在 IDL 中为 ReviewGameVersionRequest 添加一个可选的 comment 字段。
	reviewComment := ""

	// --- 3. 调用 DAO 层更新数据库 ---
	err := GameDao.ReviewGameVersion(ctx, uint64(req.GameID), uint64(req.GameVersionID), newStatus, reviewComment)
	if err != nil {
		// 如果 DAO 返回 "记录未找到" 错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &game.ReviewGameVersionResponse{
				BaseResp: &common.BaseResp{Code: "10002", Msg: "Game or Version not found"},
			}, nil
		}
		// 其他数据库错误
		return &game.ReviewGameVersionResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "Failed to update game version status: " + err.Error()},
		}, nil
	}

	// --- 4. 构建并返回成功的响应 ---
	resp := &game.ReviewGameVersionResponse{
		BaseResp: &common.BaseResp{Code: "200", Msg: "Success"},
	}
	return resp, nil
}
