package main

import (
	"log"

	game "github.com/GameLaunchPad/game_management_project/kitex_gen/game/gameservice"
)

func main() {
	svr := game.NewServer(new(GameServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
