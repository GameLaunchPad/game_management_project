package dal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GameLaunchPad/game_management_project/constdef"
	"github.com/GameLaunchPad/game_management_project/handler"
	"github.com/GameLaunchPad/game_management_project/repository"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 1. 修改函数签名，让它返回 handler 和 error
func InitClient(ctx context.Context) (*handler.CPMaterialHandler, error) {
	initIDGenerator(ctx)

	// 2. 检查 initDB 是否出错
	if err := initDB(ctx); err != nil {
		// 如果有错误，则返回 nil 和错误信息
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	cpMaterialRepo := repository.NewCPMaterialRepo(DB)
	cpMaterialHandler := handler.NewCPMaterialHandler(cpMaterialRepo)

	// 3. 在函数末尾返回创建好的实例和 nil (表示成功)
	return cpMaterialHandler, nil
}

// (推荐修改) 让 initDB 返回 error 而不是 panic
func initDB(ctx context.Context) error {
	var err error
	dsn := "kitex_user:your_actual_password@tcp(127.0.0.1:3306)/kitex_db?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果连接失败，返回错误
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		// 获取 sql.DB 失败，也返回错误
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection initialized successfully.")
	return nil // 表示成功
}

func initIDGenerator(ctx context.Context) {
	var options = idgen.NewIdGeneratorOptions(constdef.IDWorkers)
	idgen.SetIdGenerator(options)
}
