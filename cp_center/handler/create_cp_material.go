package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/yitter/idgenerator-go/idgen"
)

func (h *CPMaterialHandler) CreateCPMaterial(ctx context.Context, req *cp_center.CreateCPMaterialRequest) (*cp_center.CreateCPMaterialResponse, error) {
	// 1. 输入参数校验
	if req.CPMaterial == nil || req.CPMaterial.CpID == 0 || req.CPMaterial.CpName == "" {
		return nil, fmt.Errorf("invalid arguments: CPMaterial, CpID, and CpName are required")
	}

	materialInfo := req.GetCPMaterial()

	// 2. 根据提交模式确定状态
	var status int32
	switch req.SubmitMode {
	case cp_center.SubmitMode_SubmitDraft:
		status = int32(cp_center.MaterialStatus_Draft)
	case cp_center.SubmitMode_SubmitReview:
		status = int32(cp_center.MaterialStatus_Reviewing)
	default:
		status = int32(cp_center.MaterialStatus_Draft)
	}

	// 3. 序列化图片数组
	var verificationImagesJSON string
	if len(materialInfo.VerificationImages) > 0 {
		jsonData, err := json.Marshal(materialInfo.VerificationImages)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal verification images: %w", err)
		}
		verificationImagesJSON = string(jsonData)
	}

	// 4. 构造数据库 model 对象
	newMaterialID := idgen.NextId()

	material := &ddl.GpCpMaterial{
		Id:                 uint64(newMaterialID),
		CpId:               uint64(materialInfo.CpID),
		CpIcon:             materialInfo.CpIcon,
		CpName:             materialInfo.CpName,
		VerificationImages: verificationImagesJSON,
		BusinessLicense:    materialInfo.BusinessLicenses,
		Website:            materialInfo.Website,
		Status:             int(status),
		CreateTs:           time.Now(),
		ModifyTs:           time.Now(),
	}

	if err := h.Repo.CreateMaterial(ctx, material); err != nil {
		return &cp_center.CreateCPMaterialResponse{
			BaseResp: &common.BaseResp{
				Code: "500",
				Msg:  fmt.Sprintf("failed to create cp material in db: %s", err.Error()),
			},
		}, nil
	}

	resp := &cp_center.CreateCPMaterialResponse{
		CpID:       int64(material.CpId),
		MaterialID: int64(material.Id),
		BaseResp: &common.BaseResp{
			Code: "0",
			Msg:  "创建成功",
		},
	}

	return resp, nil
}
