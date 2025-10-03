package main

import (
	cp_center "github.com/GameLaunchPad/game_management_project/kitex_gen/cp_center/cpcenterservice"
	"log"
)

func main() {
	svr := cp_center.NewServer(new(CpCenterServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
