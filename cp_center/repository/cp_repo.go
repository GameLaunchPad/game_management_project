package repository // 或者你的包名

import (
	"context"
	"errors"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl" // 确认路径是否正确
	"gorm.io/gorm"
)

type cpRepoImpl struct {
	db *gorm.DB
}

func NewCPRepo(db *gorm.DB) ICPRepo { // 假设你的接口叫 ICPRepo
	return &cpRepoImpl{db: db}
}

// CreateCP 用于向数据库中插入一条新的厂商记录
func (c *cpRepoImpl) CreateCP(ctx context.Context, cp *ddl.GpCp) error {
	// 使用 WithContext 将上下文传递给 GORM，以便于控制超时和取消
	// Create 方法会将 cp 对象直接存入 gp_cp 表中
	result := c.db.WithContext(ctx).Create(cp)
	return result.Error
}

func (c *cpRepoImpl) GetCPByID(ctx context.Context, cpID int64) (*ddl.GpCp, error) {
	// 声明一个 ddl.GpCp 类型的变量 cp，用于存储查询结果
	var cp ddl.GpCp

	// 使用 GORM 的 First 方法来查询记录
	// .WithContext(ctx) -> 传递上下文，用于控制请求的超时和取消
	// .Where("id = ?", cpID) -> 设置查询条件，根据主键 id 查找
	// .First(&cp) -> 执行查询，并将找到的第一条记录填充到 cp 变量中
	// 如果没有找到记录，First 方法会返回 gorm.ErrRecordNotFound 错误
	result := c.db.WithContext(ctx).Where("id = ?", cpID).First(&cp)

	// 检查查询过程中是否发生错误
	if result.Error != nil {
		// 如果有错误（包括未找到记录的错误），则返回 nil 和错误信息
		return nil, result.Error
	}

	// 如果查询成功，返回找到的厂商信息指针和 nil
	return &cp, nil
}

// UpdateCP 用于更新指定 ID 的厂商信息
func (c *cpRepoImpl) UpdateCP(ctx context.Context, cpID int64, updates map[string]interface{}) error {
	// 首先，检查 updates map 是否为空，避免不必要的数据库操作
	if len(updates) == 0 {
		return errors.New("update data is empty")
	}

	// Model(&ddl.GpCp{}) -> 指定要操作的是 gp_cp 表
	// Where("id = ?", cpID) -> 添加查询条件，找到那条需要更新的记录
	// Updates(updates) -> GORM将 map 中的键值对更新到对应的列
	result := c.db.WithContext(ctx).Model(&ddl.GpCp{}).Where("id = ?", cpID).Updates(updates)

	// 检查是否有错误发生
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的有记录被更新了。
	// 如果 RowsAffected 为 0，说明没有找到 cpID 对应的记录。
	if result.RowsAffected == 0 {
		// 返回 GORM 预定义的 ErrRecordNotFound 错误，调用方可以方便地判断是否是“未找到”的错误
		return gorm.ErrRecordNotFound
	}

	return nil
}
