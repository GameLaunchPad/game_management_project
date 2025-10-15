package rpc

import (
	"log"

	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center/cpcenterservice"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/config"
	"github.com/cloudwego/kitex/client"
)

var CPCenterClient cpcenterservice.Client

func initCpCenterClient() {
	c, err := cpcenterservice.NewClient("game", client.WithHostPorts(config.Config.Rpc.GameServiceAddr))
	if err != nil {
		log.Fatal(err)
	}
	CPCenterClient = c
}
