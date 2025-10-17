package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

// CreateCPMaterial 是优化后的主函数，负责编排整个创建流程。
func (h *CPMaterialHandler) CreateCPMaterial(ctx context.Context, req *cp_center.CreateCPMaterialRequest) (*cp_center.CreateCPMaterialResponse, error) {
	// 1. 输入参数校验 (调用独立的校验函数)
	if err := validateCreateRequest(req); err != nil {
		// 返回标准的 InvalidArgument 错误，上层框架（如gRPC）可以将其转换为对应的状态码
		return nil, err // 这是关键改动：直接返回 error
	}
	log.Printf("CreateCPMaterial 参数校验成功\n")
	// 2. 从请求构建数据库模型 (调用独立的辅助函数)
	material, err := newMaterialFromRequest(req)
	if err != nil {
		// 如果在构建过程中出错（如JSON序列化失败），也直接返回 error
		return nil, fmt.Errorf("failed to build material model from request: %w", err)
	}
	log.Printf("CreateCPMaterial 构建数据库模型material成功\n")
	// 3. 执行数据库插入操作
	if err := h.MaterialRepo.CreateMaterial(ctx, material); err != nil {
		// 数据库错误是内部服务错误，包装后返回
		// 同样，直接返回 error
		return nil, fmt.Errorf("failed to create cp material in db: %w", err)
	}

	// 4. 查看是否已有该 cp
	exists, err := h.checkCPExists(ctx, int64(material.CpId))
	if err != nil {
		return nil, fmt.Errorf("failed to check if cp exists: %w", err)
	}
	if exists {
		log.Printf("CreateCPMaterial CP 已存在\n")
		// 如果 CP 存在，则更新其最新资质ID和名称
		updates := map[string]interface{}{
			"newest_material_id": material.Id,
			"cp_name":            material.CpName, // 同时更新名称以保持最新
		}
		if err := h.CPRepo.UpdateCP(ctx, int64(material.CpId), updates); err != nil {
			return nil, fmt.Errorf("failed to update existing cp: %w", err)
		}
	} else {
		log.Printf("CreateCPMaterial CP 不存在\n")
		// 如果 CP 不存在，则创建新的 CP
		cp, err := newCPFromMaterial(material)
		// newCPFromMaterial 目前不会返回 error，但保留检查以防未来修改
		if err != nil {
			return nil, fmt.Errorf("failed to build cp model from material: %w", err)
		}
		if err := h.CPRepo.CreateCP(ctx, cp); err != nil {
			return nil, fmt.Errorf("failed to create cp in db: %w", err)
		}
	}

	// 5. 构造成功响应
	// 只有在完全成功时，才返回 (response, nil)
	resp := &cp_center.CreateCPMaterialResponse{
		CpID:       int64(material.CpId),
		MaterialID: int64(material.Id),
		BaseResp: &common.BaseResp{
			Code: "0", // 成功码
			Msg:  "创建成功",
		},
	}

	return resp, nil
}

// validateCreateRequest 负责校验输入参数的合法性。
func validateCreateRequest(req *cp_center.CreateCPMaterialRequest) error {
	if req.GetCPMaterial() == nil {
		return fmt.Errorf("CPMaterial cannot be nil") // 使用更具体的错误信息
	}
	if req.CPMaterial.GetCpID() == 0 {
		return fmt.Errorf("CpID is required")
	}
	// 使用 strings.TrimSpace 避免名称只包含空格的情况
	if strings.TrimSpace(req.CPMaterial.GetCpName()) == "" {
		return fmt.Errorf("CpName is required and cannot be empty or whitespace")
	}
	// 这里还可以添加其他校验
	return nil
}

// newMaterialFromRequest 负责将请求对象转换为数据库模型。
// 这种函数也称为 "Converter" 或 "Builder"。
func newMaterialFromRequest(req *cp_center.CreateCPMaterialRequest) (*ddl.GpCpMaterial, error) {
	materialInfo := req.GetCPMaterial()

	// 确定状态
	var status cp_center.MaterialStatus
	switch req.SubmitMode {
	case cp_center.SubmitMode_SubmitReview:
		status = cp_center.MaterialStatus_Reviewing
	case cp_center.SubmitMode_SubmitDraft:
		fallthrough // fallthrough 可以让 SubmitDraft 和 default 执行同样逻辑
	default:
		status = cp_center.MaterialStatus_Draft
	}

	// 序列化图片数组
	var verificationImagesJSON string
	if len(materialInfo.VerificationImages) > 0 {
		jsonData, err := json.Marshal(materialInfo.VerificationImages)
		if err != nil {
			// 如果序列化失败，这是一个明确的错误，需要返回
			return nil, fmt.Errorf("failed to marshal verification images: %w", err)
		}
		verificationImagesJSON = string(jsonData)
	} else {
		// 明确地给一个空的JSON数组字符串，而不是空字符串，这样数据库中的数据更一致
		verificationImagesJSON = "[]"
	}

	// 生成唯一ID并构建模型
	material := &ddl.GpCpMaterial{
		Id:                 uint64(idgen.NextId()),
		CpId:               uint64(materialInfo.CpID),
		CpIcon:             materialInfo.CpIcon,
		CpName:             materialInfo.CpName,
		VerificationImages: verificationImagesJSON,
		BusinessLicense:    materialInfo.BusinessLicenses,
		Website:            materialInfo.Website,
		Status:             int(status),
		// CreateTs 和 ModifyTs 由数据库的 DEFAULT CURRENT_TIMESTAMP 和 ON UPDATE 自动处理会更好。
		// 如果数据库表结构没有设置自动处理，那么在Go中设置是必要的。
		// 假设表结构如你所给，Go中设置是正确的。
		CreateTs: time.Now(),
		ModifyTs: time.Now(),
	}

	return material, nil
}

func newCPFromMaterial(material *ddl.GpCpMaterial) (*ddl.GpCp, error) {

	cp := &ddl.GpCp{
		Id:               material.CpId,
		CpName:           material.CpName,
		NewestMaterialId: material.Id,
		OnlineMaterialId: 0,
		VerifyStatus:     0,
		CreateTs:         time.Now(),
		ModifyTs:         time.Now(),
	}
	return cp, nil
}

func (h *CPMaterialHandler) checkCPExists(ctx context.Context, id int64) (bool, error) {
	// 调用仓库的 GetCPByID 方法
	_, err := h.CPRepo.GetCPByID(ctx, id)

	// 分析从仓库返回的错误
	if err != nil {
		// 如果错误是 gorm.ErrRecordNotFound，这意味着CP不存在。
		// 这是一个预期的结果，而不是一个系统故障。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回 false (不存在) 并且没有错误。
			return false, nil
		}
		// 对于任何其他类型的错误（如数据库连接问题等），
		// 这是一个真正的问题。返回 false 并传播该错误。
		return false, err
	}

	// 如果 err 为 nil，意味着成功找到了记录。
	// 返回 true (存在) 并且没有错误。
	return true, nil
}
