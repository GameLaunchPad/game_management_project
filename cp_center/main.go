package main

import (
	"context"
	"log"

	"github.com/GameLaunchPad/game_management_project/cp_center/dal"
	cp_center "github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center/cpcenterservice"
)

func main() {
	// 1. 调用初始化函数，并接收返回的 handler
	cpMaterialHandler, err := dal.InitClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to init client: %v", err)
	}

	// 2. 创建服务实现时，将 handler "注入" 进去
	//    NewCpCenterServiceImpl 是我们需要创建的一个新函数
	serviceImpl := NewCpCenterServiceImpl(cpMaterialHandler)

	svr := cp_center.NewServer(serviceImpl)

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
