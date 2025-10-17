package handler

import (
	"context"
	"log"

	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
)

func (h *CPMaterialHandler) GetCPMaterial(ctx context.Context, req *cp_center.GetCPMaterialRequest) (*cp_center.GetCPMaterialResponse, error) {
	log.Printf("GetCPMaterial 收到请求参数：%+v\n", req)
	// 参数校验
	if req.CpID <= 0 {
		return &cp_center.GetCPMaterialResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: "invalid parameter: cp_id must be positive: "},
		}, nil
	}
	log.Printf("参数校验成功\n")
	material, err := h.MaterialRepo.GetMaterialByID(ctx, req.MaterialID)
	log.Printf("获得结果：%+v。错误：%+v\n", material, err)
	// 错误处理
	if err != nil {
		return &cp_center.GetCPMaterialResponse{
			BaseResp: &common.BaseResp{Code: "500", Msg: err.Error()},
		}, nil
	}
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
		BaseResp: &common.BaseResp{Code: "0", Msg: "success"},
	}
	log.Printf("resp: %+v\n", material)
	return resp, nil
}
