package game_platform_api

import (
	"context"

	"github.com/GameLaunchPad/game_management_project/biz/model/game_platform_api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UpdateGameDetail .
// @router /api/v1/games/:id [PUT]
func UpdateGameDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req game_platform_api.UpdateGameDetailRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(game_platform_api.UpdateGameDetailResponse)

	c.JSON(consts.StatusOK, resp)
}
