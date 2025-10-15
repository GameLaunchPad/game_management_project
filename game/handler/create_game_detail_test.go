package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/GameLaunchPad/game_management_project/game/dao/mock"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestCreateGameDetail_Success tests the successful creation of a new game detail
func TestCreateGameDetail_Success(t *testing.T) {
	setupIDGenerator()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	// expect CreateGame to be called once with any parameters and return nil error
	mockGameDAO.EXPECT().CreateGame(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	req := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 0,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "My First Game",
			},
		},
	}

	resp, err := CreateGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
	assert.NotEqual(t, int64(0), resp.GameID)
}

// TestCreateGameDetail_FailWithNonZeroGameID tests the failure case when a non-zero GameID is provided
func TestCreateGameDetail_FailWithNonZeroGameID(t *testing.T) {
	req := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 123,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "A Game That Should Not Be Created",
			},
		},
	}

	resp, err := CreateGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "GameID must be 0")
}

// TestCreateGameDetail_DaoError tests the failure case when the DAO returns an error
func TestCreateGameDetail_DaoError(t *testing.T) {
	setupIDGenerator()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	mockedError := errors.New("a generic database error")
	mockGameDAO.EXPECT().CreateGame(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockedError).Times(1)

	req := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 0,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "Another Game",
			},
		},
	}

	resp, err := CreateGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "500", resp.BaseResp.Code)
}
