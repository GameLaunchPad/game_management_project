package handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"gorm.io/gorm"
)

func (h *CPMaterialHandler) UpdateCPMaterial(ctx context.Context, req *cp_center.UpdateCPMaterialRequest) (*cp_center.UpdateCPMaterialResponse, error) {
	// 参数校验
	if req.MaterialID <= 0 {
		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  "invalid parameter: material_id is required",
			},
		}, nil
	}
	if req.CpMaterial == nil {
		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  "invalid parameter: cp_material data is missing",
			},
		}, nil
	}
	if req.SubmitMode == cp_center.SubmitMode_Unset {
		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  "invalid parameter: submit_mode is required",
			},
		}, nil
	}
	if req.CpMaterial.CpName == "" || req.CpMaterial.BusinessLicenses == "" {
		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  "cp_name and business_license are required fields",
			},
		}, nil
	}

	// 查询原始记录
	material, err := h.MaterialRepo.GetMaterialByID(ctx, req.MaterialID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &cp_center.UpdateCPMaterialResponse{
				BaseResp: &common.BaseResp{
					Code: "500",
					Msg:  "material not found",
				},
			}, nil
		}
		return nil, err
	}

	// 业务逻辑校验
	if material.Status == 2 || material.Status == 3 { // 2-审核中, 3-已发布
		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  "cannot update material that is in review or online",
			},
		}, nil
	}

	// 准备要更新的数据
	updates := make(map[string]interface{})
	updates["cp_icon"] = req.CpMaterial.CpIcon
	updates["cp_name"] = req.CpMaterial.CpName
	updates["business_license"] = req.CpMaterial.BusinessLicenses
	updates["website"] = req.CpMaterial.Website

	if req.CpMaterial.VerificationImages != nil {
		imgBytes, errJson := json.Marshal(req.CpMaterial.VerificationImages)
		if errJson != nil {
			return nil, errors.New("failed to marshal verification_images")
		}
		updates["verification_images"] = string(imgBytes)
	}

	switch req.SubmitMode {
	case cp_center.SubmitMode_SubmitDraft:
		updates["status"] = 1 // 1-草稿
	case cp_center.SubmitMode_SubmitReview:
		updates["status"] = 2 // 2-审核中
	}

	updates["modify_ts"] = time.Now()

	// 执行更新操作
	_, err = h.MaterialRepo.UpdateMaterial(ctx, req.MaterialID, updates)
	if err != nil {

		return &cp_center.UpdateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  err.Error(),
			},
		}, nil
	}

	// 构建并返回成功响应
	resp := &cp_center.UpdateCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: "0",
			Msg:  "success",
		},
	}
	return resp, nil
}
