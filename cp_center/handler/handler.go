package handler

import "github.com/GameLaunchPad/game_management_project/cp_center/repository"

type CPMaterialHandler struct {
	MaterialRepo repository.ICPMaterialRepo
	CPRepo       repository.ICPRepo
}

// NewCPMaterialHandler 是 Handler 的构造函数
func NewCPMaterialHandler(MaterialRepo repository.ICPMaterialRepo, CPRepo repository.ICPRepo) *CPMaterialHandler {
	return &CPMaterialHandler{MaterialRepo: MaterialRepo, CPRepo: CPRepo}
}
