package main

import (
	"context"
	"log"

	"github.com/GameLaunchPad/game_management_project/dal"
	cp_center "github.com/GameLaunchPad/game_management_project/kitex_gen/cp_center/cpcenterservice"
)

func main() {
	dal.InitClient(context.Background())
	svr := cp_center.NewServer(new(CpCenterServiceImpl))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
