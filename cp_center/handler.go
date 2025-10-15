package main

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
)

// CpCenterServiceImpl implements the last service interface defined in the IDL.
type CpCenterServiceImpl struct {
	CpMaterialHandler *handler.CPMaterialHandler
}

func NewCpCenterServiceImpl(handler *handler.CPMaterialHandler) *CpCenterServiceImpl {
	return &CpCenterServiceImpl{
		CpMaterialHandler: handler,
	}
}

// CreateCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) CreateCPMaterial(ctx context.Context, req *cp_center.CreateCPMaterialRequest) (resp *cp_center.CreateCPMaterialResponse, err error) {
	return s.CpMaterialHandler.CreateCPMaterial(ctx, req)
}

// UpdateCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) UpdateCPMaterial(ctx context.Context, req *cp_center.UpdateCPMaterialRequest) (resp *cp_center.UpdateCPMaterialResponse, err error) {
	return s.CpMaterialHandler.UpdateCPMaterial(ctx, req)
}

// ReviewCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) ReviewCPMaterial(ctx context.Context, req *cp_center.ReviewCPMaterialRequest) (resp *cp_center.ReviewCPMaterialResponse, err error) {
	return s.CpMaterialHandler.ReviewCPMaterial(ctx, req)
}

// GetCPMaterial implements the CpCenterServiceImpl interface.
func (s *CpCenterServiceImpl) GetCPMaterial(ctx context.Context, req *cp_center.GetCPMaterialRequest) (resp *cp_center.GetCPMaterialResponse, err error) {
	return s.CpMaterialHandler.GetCPMaterial(ctx, req)
}
