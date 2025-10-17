// Package repository 定义了数据访问层的接口。
// 这个文件是 mockgen 工具生成 mock 代码的“图纸”，其内容应与 cp_material_repo.go 中的接口定义保持完全一致。
package repository

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
)

type ICPMaterialRepo interface {
	// CreateMaterial 用于创建新的 CP 素材记录
	CreateMaterial(ctx context.Context, material *ddl.GpCpMaterial) error

	// UpdateMaterial 用于根据素材 ID 更新指定的字段
	UpdateMaterial(ctx context.Context, materialID int64, updates map[string]interface{}) (int64, error)

	// GetMaterialByID 用于根据素材 ID 获取单个素材的详细信息
	GetMaterialByID(ctx context.Context, materialID int64) (*ddl.GpCpMaterial, error)

	// GetMaterialByCPID 用于根据 CP ID 获取单个素材的详细信息
	GetMaterialByCPID(ctx context.Context, cpID int64) (*ddl.GpCpMaterial, error)
}

type ICPRepo interface {
	CreateCP(ctx context.Context, cp *ddl.GpCp) error
	GetCPByID(ctx context.Context, cpID int64) (*ddl.GpCp, error)
	UpdateCP(ctx context.Context, cpID int64, updates map[string]interface{}) error
}
