package rpc

import (
	"log"

	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game/gameservice"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/config"
	"github.com/cloudwego/kitex/client"
)

var GameClient gameservice.Client

func initGameClient() {
	c, err := gameservice.NewClient("game", client.WithHostPorts(config.Config.Rpc.GameServiceAddr))
	if err != nil {
		log.Fatal(err)
	}
	GameClient = c
}
