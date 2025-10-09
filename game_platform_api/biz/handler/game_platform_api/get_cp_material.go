package game_platform_api

import (
	"github.com/GameLaunchPad/game_management_project/biz/model/game_platform_api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetCPMaterial .
// @router /api/v1/cp/materials/:id [GET]
func GetCPMaterial(ctx context.Context, c *app.RequestContext) {
	var err error
	var req game_platform_api.GetCPMaterialRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(game_platform_api.GetCPMaterialResponse)

	c.JSON(consts.StatusOK, resp)
}
