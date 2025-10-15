package handler

import "github.com/GameLaunchPad/game_management_project/repository"

type CPMaterialHandler struct {
	Repo repository.ICPMaterialRepo
}

// NewCPMaterialHandler 是 Handler 的构造函数
func NewCPMaterialHandler(repo repository.ICPMaterialRepo) *CPMaterialHandler {
	return &CPMaterialHandler{Repo: repo}
}
