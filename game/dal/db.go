package dal

import (
	"log"

	"github.com/GameLaunchPad/game_management_project/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	if config.GlobalConfig == nil {
		panic("config not initialized")
	}

	dsn := config.GlobalConfig.MySQL.DSN
	if dsn == "" {
		panic("MySQL DSN is empty")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	log.Println("Connected to database successfully")
}
