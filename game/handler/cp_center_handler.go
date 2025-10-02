package handler

import (
	"context"
	"strconv"

	"github.com/GameLaunchPad/game_management_project/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/cp_center"
)

// CpCenterServiceImpl 实现了 CpCenterService 接口
// 同样，如果需要数据库等依赖，通过字段注入进来
type CpCenterServiceImpl struct {
	// DbClient *gorm.DB // 举例：注入数据库客户端
}

// NewCpCenterServiceImpl 是 CpCenterServiceImpl 的构造函数
func NewCpCenterServiceImpl( /* db *gorm.DB */ ) *CpCenterServiceImpl {
	return &CpCenterServiceImpl{ /* DbClient: db */ }
}

// ----------- 在下面实现 thrift 中定义的每一个方法 -----------

// CreateCPMaterial 实现了创建认证材料的逻辑
func (s *CpCenterServiceImpl) CreateCPMaterial(ctx context.Context, req *cp_center.CreateCPMaterialRequest) (resp *cp_center.CreateCPMaterialResponse, err error) {
	// TODO: 在这里实现你的业务逻辑
	// 1. 参数校验，检查 req.CPMaterial 是否合法。
	// 2. 将 req.CPMaterial (thrift 生成的 struct) 转换为你的数据库模型 (model)。
	// 3. 根据 req.SubmitMode 设置材料的状态 (是草稿还是待审核)。
	// 4. 调用 DAL/DAO 层将数据写入数据库。
	// 5. 组装响应 CreateCPMaterialResponse 并返回。

	resp = &cp_center.CreateCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: strconv.Itoa(0),
			Msg:  "创建成功",
		},
		CpID:       1001, // 示例ID
		MaterialID: 2001, // 示例ID
	}
	return resp, nil
}

// UpdateCPMaterial 实现了更新认证材料的逻辑
func (s *CpCenterServiceImpl) UpdateCPMaterial(ctx context.Context, req *cp_center.UpdateCPMaterialRequest) (resp *cp_center.UpdateCPMaterialResponse, err error) {
	// TODO: 在这里实现你的业务逻辑
	// 1. 根据 req.MaterialID 检查材料是否存在。
	// 2. 更新数据库中的材料信息。
	// 3. 根据 req.SubmitMode 更新材料的状态。
	// 4. 返回操作结果。

	resp = &cp_center.UpdateCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: strconv.Itoa(0),
			Msg:  "update成功",
		},
	}
	return resp, nil
}

// ReviewCPMaterial 实现了审核厂商材料的逻辑
func (s *CpCenterServiceImpl) ReviewCPMaterial(ctx context.Context, req *cp_center.ReviewCPMaterialRequest) (resp *cp_center.ReviewCPMaterialResponse, err error) {
	// TODO: 在这里实现你的业务逻辑
	// 1. 根据 req.MaterialID 找到对应的材料。
	// 2. 根据 req.ReviewResult (通过或拒绝) 更新材料的状态。
	// 3. 如果是拒绝，可能需要记录拒绝原因。
	// 4. 返回操作结果。

	resp = &cp_center.ReviewCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: strconv.Itoa(0),
			Msg:  "review成功",
		},
	}
	return resp, nil
}

// GetCPMaterial 实现了获取厂商认证材料的逻辑
func (s *CpCenterServiceImpl) GetCPMaterial(ctx context.Context, req *cp_center.GetCPMaterialRequest) (resp *cp_center.GetCPMaterialResponse, err error) {
	// TODO: 在这里实现你的业务逻辑
	// 1. 根据 req.CpID 从数据库中查询材料信息。
	// 2. 如果找不到，返回相应的错误信息。
	// 3. 将从数据库查出的数据模型 (model) 转换为 thrift 定义的 CPMaterial 结构体。
	// 4. 组装响应并返回。

	// 示例返回
	resp = &cp_center.GetCPMaterialResponse{
		CPMaterial: &cp_center.CPMaterial{
			MaterialID: 2001,
			CpID:       req.CpID,
			CpName:     "示例厂商",
			Status:     cp_center.MaterialStatus_Online,
		},
		BaseResp: &common.BaseResp{
			Code: strconv.Itoa(0),
			Msg:  "get成功",
		},
	}
	return resp, nil
}
