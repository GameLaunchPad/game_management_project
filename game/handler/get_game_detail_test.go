package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/game/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/game/dao/mock"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetGameDetail_Success tests the successful retrieval of game details.
func TestGetGameDetail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// 1. prepare mock data
	gameID := uint64(123)
	mockGame := &ddl.GpGame{Id: gameID, CpId: 1001, GameName: "My Detail Test Game", CreateTs: time.Now(), ModifyTs: time.Now()}
	mockNewestVersion := &ddl.GpGameVersion{Id: 201, GameId: gameID, GameName: "Version 1.1", Platform: "[]", GameIntroductionImages: "[]"}
	mockOnlineVersion := &ddl.GpGameVersion{Id: 200, GameId: gameID, GameName: "Version 1.0", Platform: "[]", GameIntroductionImages: "[]"}

	// 2. define expectations
	mockGameDAO.EXPECT().
		GetGameDetail(gomock.Any(), gameID).
		Return(mockGame, mockNewestVersion, mockOnlineVersion, nil).
		Times(1)

	// 3. prepare request
	req := &game.GetGameDetailRequest{
		GameID: int64(gameID),
	}

	// 4. call the function
	resp, err := GetGameDetail(context.Background(), req)

	// 5. assert results
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
	assert.NotNil(t, resp.GameDetail)
	assert.Equal(t, int64(gameID), resp.GameDetail.GameID)
	assert.Equal(t, "Version 1.1", resp.GameDetail.NewestGameVersion_.GameName)
	assert.Equal(t, "Version 1.0", resp.GameDetail.OnlineGameVersion.GameName)
}

// TestGetGameDetail_GameNotFound tests the scenario where the requested game does not exist.
func TestGetGameDetail_GameNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(999) // invalid/non-existent game ID

	// define expectations: DAO returns gorm.ErrRecordNotFound
	mockGameDAO.EXPECT().
		GetGameDetail(gomock.Any(), gameID).
		Return(nil, nil, nil, gorm.ErrRecordNotFound).
		Times(1)

	req := &game.GetGameDetailRequest{
		GameID: int64(gameID),
	}

	resp, err := GetGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10001", resp.BaseResp.Code, "Expected game not found error code")
	assert.Equal(t, "Game not found", resp.BaseResp.Msg)
	assert.Nil(t, resp.GameDetail, "GameDetail should be nil on error")
}

// TestGetGameDetail_DaoError tests the scenario where the DAO returns a general error.
func TestGetGameDetail_DaoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(456)
	mockedError := errors.New("unexpected database error")

	// define expectations: DAO returns a general error
	mockGameDAO.EXPECT().
		GetGameDetail(gomock.Any(), gameID).
		Return(nil, nil, nil, mockedError).
		Times(1)

	req := &game.GetGameDetailRequest{
		GameID: int64(gameID),
	}

	resp, err := GetGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "500", resp.BaseResp.Code, "Expected internal server error code")
}
