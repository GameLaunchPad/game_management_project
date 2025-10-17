package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/biz/model/game_platform_api"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/rpc"
)

type CpCenterService struct{}

func NewCpCenterService() *CpCenterService {
	return &CpCenterService{}
}

// CreateMaterial 处理创建新 CP 物料的请求。
func (s *CpCenterService) CreateMaterial(ctx context.Context, req *game_platform_api.CreateCPMaterialsRequest) (*cp_center.CreateCPMaterialResponse, error) {
	log.Printf("CreateMaterial 收到请求参数：%+v\n", req)
	if req.CpMaterial == nil {
		return nil, fmt.Errorf("请求中的 CpMaterial 为空")
	}

	cpID, err := strconv.ParseInt(req.CpMaterial.CpID, 10, 64)
	if err != nil {
		// --- 修改点 1: 包装错误，提供更清晰的错误链 ---
		return nil, fmt.Errorf("无效的 CpID 格式：%w", err)
	}

	// 构造 RPC 请求逻辑完全不变
	rpcReq := &cp_center.CreateCPMaterialRequest{
		CPMaterial: &cp_center.CPMaterial{
			CpID:               cpID,
			CpName:             req.CpMaterial.CpName,
			CpIcon:             req.CpMaterial.CpIcon,
			Website:            req.CpMaterial.Website,
			VerificationImages: req.CpMaterial.VerificationImages,
			BusinessLicenses:   req.CpMaterial.BusinessLicense,
		},
		SubmitMode: cp_center.SubmitMode(req.SubmitMode),
	}

	// 保持您原有的全局客户端调用
	resp, err := rpc.CPCenterClient.CreateCPMaterial(ctx, rpcReq)
	if err != nil {
		// --- 修改点 2: 同样包装错误 ---
		return nil, fmt.Errorf("RPC 调用 CreateCPMaterial 失败：%w", err)
	}

	// --- 修改点 3: 将业务错误也转为标准 error 返回 ---
	if resp.BaseResp != nil && resp.BaseResp.GetCode() != "0" {
		// 原来您这里是直接返回 resp 和 nil，现在统一返回 nil 和一个描述业务错误的 error
		return nil, fmt.Errorf("业务错误：%s", resp.BaseResp.Msg)
	}
	// --- 修改结束 ---

	return resp, nil
}

// GetMaterial 处理获取 CP 物料信息的请求
func (s *CpCenterService) GetMaterial(ctx context.Context, req *game_platform_api.GetCPMaterialRequest) (*cp_center.GetCPMaterialResponse, error) {
	log.Printf("GetMaterial 收到请求参数：%+v\n", req)
	cpID, err := strconv.ParseInt(req.CpID, 10, 64)
	materialID, err := strconv.ParseInt(req.MaterialID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的 CpID 格式，必须为数字字符串：%w", err)
	}

	// 构造 RPC 请求
	rpcReq := &cp_center.GetCPMaterialRequest{
		CpID:       cpID,
		MaterialID: materialID,
	}

	// 调用下游 RPC 服务
	resp, err := rpc.CPCenterClient.GetCPMaterial(ctx, rpcReq)
	log.Printf("resp：%+v, err: %+v\n", resp, err)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// 业务逻辑错误处理
	if resp.BaseResp != nil && resp.BaseResp.GetCode() != "0" {
		return nil, fmt.Errorf("业务错误：%s", resp.BaseResp.Msg)
	}
	return resp, nil
}

// ReviewMaterial 处理审核 CP 物料的请求。
func (s *CpCenterService) ReviewMaterial(ctx context.Context, req *game_platform_api.ReviewCPMaterialRequest) (*cp_center.ReviewCPMaterialResponse, error) {
	log.Printf("ReviewMaterial 收到请求参数：%+v\n", req)
	// 构造 RPC 请求
	rpcReq := &cp_center.ReviewCPMaterialRequest{
		CpID:          req.CpID,
		MaterialID:    req.MaterialID,
		ReviewResult_: cp_center.ReviewResult_(req.ReviewResult),
		ReviewRemark: &cp_center.ReviewRemark{
			Remark:     req.ReviewRemark.Remark,
			Operator:   req.ReviewRemark.Operator,
			ReviewTime: req.ReviewRemark.ReviewTime,
			Meta:       req.ReviewRemark.Meta,
		},
	}

	// 调用下游 RPC 服务
	resp, err := rpc.CPCenterClient.ReviewCPMaterial(ctx, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("RPC 调用 ReviewCPMaterial 失败：%w", err)
	}

	// 业务逻辑错误处理
	if resp.BaseResp != nil && resp.BaseResp.GetCode() != "0" {
		return nil, fmt.Errorf("业务错误：%s", resp.BaseResp.Msg)
	}

	return resp, nil
}

// UpdateMaterial 处理更新已存在 CP 物料的请求
func (s *CpCenterService) UpdateMaterial(ctx context.Context, req *game_platform_api.UpdateCPMaterialsRequest) (*cp_center.UpdateCPMaterialResponse, error) {
	log.Printf("UpdateMaterial 收到请求参数：%+v\n", req)
	if req.CpMaterial == nil {
		return nil, fmt.Errorf("请求中的 CpMaterial 为空")
	}

	// 类型转换：将 API 层的 string 类型 cp_id 转换为 RPC 层的 int64
	cpID, err := strconv.ParseInt(req.CpMaterial.CpID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的 CpID 格式：%w", err)
	}

	// 构造 RPC 请求
	rpcReq := &cp_center.UpdateCPMaterialRequest{
		MaterialID: req.MaterialID,
		CpMaterial: &cp_center.CPMaterial{
			CpID:               cpID,
			CpName:             req.CpMaterial.CpName,
			CpIcon:             req.CpMaterial.CpIcon,
			Website:            req.CpMaterial.Website,
			VerificationImages: req.CpMaterial.VerificationImages,
			BusinessLicenses:   req.CpMaterial.BusinessLicense,
		},
		SubmitMode: cp_center.SubmitMode(req.SubmitMode),
	}

	// 调用下游 RPC 服务
	resp, err := rpc.CPCenterClient.UpdateCPMaterial(ctx, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("RPC 调用 UpdateCPMaterial 失败：%w", err)
	}

	// 业务逻辑错误处理
	if resp.BaseResp != nil && resp.BaseResp.GetCode() != "0" {
		return nil, fmt.Errorf("业务错误：%s", resp.BaseResp.Msg)
	}

	return resp, nil
}
