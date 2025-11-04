// file: game/handler/get_game_list_test.go
package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/game/dao"
	"github.com/GameLaunchPad/game_management_project/game/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/game/dao/mock"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestGetGameList_Success 测试获取游戏列表的成功场景
func TestGetGameList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 1. 准备 mock DAO 的返回值
	mockedGameList := []*dao.GameWithVersionStatus{
		{
			GpGame: ddl.GpGame{
				Id:       1,
				CpId:     1001,
				GameName: "Test Game 1",
				CreateTs: time.Now(),
				ModifyTs: time.Now(),
			},
			Status: int(game.GameStatus_Published),
		},
		{
			GpGame: ddl.GpGame{
				Id:       2,
				CpId:     1002,
				GameName: "Test Game 2",
				CreateTs: time.Now(),
				ModifyTs: time.Now(),
			},
			Status: int(game.GameStatus_Draft),
		},
	}
	var mockTotal int64 = 2

	// 2. 定义期望：GetGameList 方法被调用1次，并返回我们准备好的数据
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), nil, 1, 10). // 期望在没有过滤器、第一页、每页10条的情况下被调用
		Return(mockedGameList, mockTotal, nil).
		Times(1)

	// 3. 准备请求参数
	req := &game.GetGameListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 4. 调用被测试的函数
	resp, err := GetGameList(context.Background(), req)

	// 5. 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
	assert.Equal(t, int32(mockTotal), resp.TotalCount)
	assert.Len(t, resp.GameList, 2)
	assert.Equal(t, "Test Game 1", resp.GameList[0].GameName)
	assert.Equal(t, game.GameStatus_Draft, resp.GameList[1].GameStatus) // 验证 Status 也被正确转换
}

// TestGetGameList_WithFilter 测试带过滤条件查询的成功场景
func TestGetGameList_WithFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	filterText := "Filtered"
	mockedGameList := []*dao.GameWithVersionStatus{
		{
			GpGame: ddl.GpGame{Id: 3, GameName: "Filtered Game"},
			Status: int(game.GameStatus_Reviewing),
		},
	}
	var mockTotal int64 = 1

	// 定义期望：这次我们期望 GetGameList 的第二个参数（filterText）不再是 nil
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), &filterText, 1, 10).
		Return(mockedGameList, mockTotal, nil).
		Times(1)

	req := &game.GetGameListRequest{
		Filter: &game.GameListFilter{
			FilterText: &filterText,
		},
		PageNum:  1,
		PageSize: 10,
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
	assert.Equal(t, int32(1), resp.TotalCount)
	assert.Len(t, resp.GameList, 1)
	assert.Equal(t, "Filtered Game", resp.GameList[0].GameName)
}

// TestGetGameList_DaoError 测试 DAO 层返回错误的场景
func TestGetGameList_DaoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 定义期望：GetGameList 方法被调用，但返回一个模拟的数据库错误
	mockedError := errors.New("database connection lost")
	mockGameDAO.EXPECT().GetGameList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), mockedError).Times(1)

	req := &game.GetGameListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "500", resp.BaseResp.Code) // 期望返回内部服务器错误码
}

// TestGetGameList_InvalidPageNum 测试无效的 PageNum（应该使用默认值 1）
func TestGetGameList_InvalidPageNum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 期望使用默认值 1
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), gomock.Any(), 1, 10).
		Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
		Times(1)

	req := &game.GetGameListRequest{
		PageNum:  0, // 无效值，应该使用默认值 1
		PageSize: 10,
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestGetGameList_NegativePageNum 测试负数的 PageNum（应该使用默认值 1）
func TestGetGameList_NegativePageNum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 期望使用默认值 1
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), gomock.Any(), 1, 10).
		Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
		Times(1)

	req := &game.GetGameListRequest{
		PageNum:  -1, // 负数，应该使用默认值 1
		PageSize: 10,
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestGetGameList_InvalidPageSize 测试无效的 PageSize（应该使用默认值 10）
func TestGetGameList_InvalidPageSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 期望使用默认值 10
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), gomock.Any(), 1, 10).
		Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
		Times(1)

	req := &game.GetGameListRequest{
		PageNum:  1,
		PageSize: 0, // 无效值，应该使用默认值 10
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestGetGameList_NegativePageSize 测试负数的 PageSize（应该使用默认值 10）
func TestGetGameList_NegativePageSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 期望使用默认值 10
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), gomock.Any(), 1, 10).
		Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
		Times(1)

	req := &game.GetGameListRequest{
		PageNum:  1,
		PageSize: -1, // 负数，应该使用默认值 10
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestGetGameList_NoFilter 测试没有 Filter 的情况（应该使用 nil）
func TestGetGameList_NoFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 期望 filterText 为 nil
	mockGameDAO.EXPECT().
		GetGameList(gomock.Any(), nil, 1, 10).
		Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
		Times(1)

	req := &game.GetGameListRequest{
		PageNum:  1,
		PageSize: 10,
		// Filter 未设置
	}

	resp, err := GetGameList(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}
