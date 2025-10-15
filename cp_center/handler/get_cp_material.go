package handler

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/cp_center"
	"gorm.io/gorm"
)

func (h *CPMaterialHandler) GetCPMaterial(ctx context.Context, req *cp_center.GetCPMaterialRequest) (*cp_center.GetCPMaterialResponse, error) {
	// 参数校验
	if req.CpID <= 0 {
		return nil, errors.New("invalid parameter: cp_id must be positive")
	}

	// 数据库查询 -> 改为调用 Repo 的方法
	material, err := h.Repo.GetMaterialByCPID(ctx, int32(req.CpID))

	// 错误处理
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cp material not found")
		}
		return nil, err
	}

	// 构建响应
	resp := &cp_center.GetCPMaterialResponse{
		CPMaterial: &cp_center.CPMaterial{
			MaterialID:         int64(material.Id),
			CpID:               int64(material.CpId),
			CpIcon:             material.CpIcon,
			CpName:             material.CpName,
			VerificationImages: nil,
			BusinessLicenses:   material.BusinessLicense,
			Website:            material.Website,
			Status:             cp_center.MaterialStatus(material.Status),
			ReviewComment:      material.ReviewComment,
			CreateTime:         material.CreateTs.Unix(),
			ModifyTime:         material.ModifyTs.Unix(),
		},
	}

	// 返回成功的响应和 nil 错误
	return resp, nil
}
