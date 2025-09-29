package main

import (
	"context"
	"log"

	"github.com/GameLaunchPad/game_management_project/dal"
	game "github.com/GameLaunchPad/game_management_project/kitex_gen/game/gameservice"
)

func main() {
	dal.InitClient(context.Background())
	svr := game.NewServer(new(GameServiceImpl))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
