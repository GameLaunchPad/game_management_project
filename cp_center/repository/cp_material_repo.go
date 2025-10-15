package repository

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"gorm.io/gorm"
)

// 注意：我们已经将 ICPMaterialRepo 接口的定义移至 interfaces.go 文件中。
// 这个文件现在只负责接口的具体实现。

// cpMaterialRepoImpl 是 ICPMaterialRepo 的 GORM 实现
type cpMaterialRepoImpl struct {
	db *gorm.DB
}

// NewCPMaterialRepo 是 cpMaterialRepoImpl 的构造函数
// 它返回的是 ICPMaterialRepo 接口类型，这样我们就将实现与外部调用者解耦了。
func NewCPMaterialRepo(db *gorm.DB) ICPMaterialRepo {
	return &cpMaterialRepoImpl{db: db}
}

// GetMaterialByCPID 实现了接口中定义的方法
func (r *cpMaterialRepoImpl) GetMaterialByCPID(ctx context.Context, cpID int32) (*ddl.GpCpMaterial, error) {
	var material ddl.GpCpMaterial
	err := r.db.WithContext(ctx).Where("cp_id = ?", cpID).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// GetMaterialByID 实现了接口中定义的方法
func (r *cpMaterialRepoImpl) GetMaterialByID(ctx context.Context, materialID int64) (*ddl.GpCpMaterial, error) {
	var material ddl.GpCpMaterial
	err := r.db.WithContext(ctx).Where("id = ?", materialID).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// UpdateMaterial 实现了接口中定义的方法
func (r *cpMaterialRepoImpl) UpdateMaterial(ctx context.Context, materialID int64, updates map[string]interface{}) (int64, error) {
	tx := r.db.WithContext(ctx).Model(&ddl.GpCpMaterial{}).Where("id = ?", materialID).Updates(updates)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// CreateMaterial 实现了接口中定义的方法
func (r *cpMaterialRepoImpl) CreateMaterial(ctx context.Context, material *ddl.GpCpMaterial) error {
	return r.db.WithContext(ctx).Create(material).Error
}
