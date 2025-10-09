package game_platform_api

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/biz/model/game_platform_api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetGameDetail .
// @router /api/v1/games/:id [GET]
func GetGameDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req game_platform_api.GetGameDetailRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(game_platform_api.GetGameDetailResponse)

	c.JSON(consts.StatusOK, resp)
}
