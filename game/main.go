package main

import (
	"context"
	"log"

	"github.com/GameLaunchPad/game_management_project/game/config"
	"github.com/GameLaunchPad/game_management_project/game/dal"
	"github.com/GameLaunchPad/game_management_project/game/dao"
	"github.com/GameLaunchPad/game_management_project/game/handler"
	game "github.com/GameLaunchPad/game_management_project/game/kitex_gen/game/gameservice"
)

const configPath = "script/config.yaml"

func main() {
	if err := config.Init(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dal.InitClient(context.Background())
	handler.GameDao = dao.NewGameDAO()

	svr := game.NewServer(new(GameServiceImpl))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
