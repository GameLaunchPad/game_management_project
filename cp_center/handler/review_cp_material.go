package handler

import (
	"context"
	"errors"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"gorm.io/gorm"
)

func (h *CPMaterialHandler) ReviewCPMaterial(ctx context.Context, req *cp_center.ReviewCPMaterialRequest) (*cp_center.ReviewCPMaterialResponse, error) {
	// 参数校验
	if req.MaterialID <= 0 {

		return nil, errors.New("invalid parameter: material_id is required")
	}
	if req.ReviewResult_ == cp_center.ReviewResult__Unset {
		return nil, errors.New("invalid parameter: review_result must be Pass or Reject")
	}

	// 查询原始记录
	_, err := h.MaterialRepo.GetMaterialByID(ctx, req.MaterialID) // 这里我们只关心是否存在，所以暂时不用 material 变量
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("material not found")
		}
		return nil, err
	}

	// 准备要更新的数据
	updates := make(map[string]interface{})
	if req.ReviewResult_ == cp_center.ReviewResult__Pass {
		updates["status"] = 3
	} else {
		updates["status"] = 4
	}
	updates["review_comment"] = req.ReviewRemark.Remark
	updates["operator"] = req.ReviewRemark.Operator
	updates["modify_ts"] = time.Now()

	// 执行更新操作
	rowsAffected, err := h.MaterialRepo.UpdateMaterial(ctx, req.MaterialID, updates)
	if err != nil {
		return nil, err // 更新失败
	}

	if rowsAffected == 0 {
		return nil, errors.New("update failed, zero rows affected")
	}

	// 构建并返回成功响应
	resp := &cp_center.ReviewCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: "0",
			Msg:  "success",
		},
	}

	return resp, nil
}
