package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/GameLaunchPad/game_management_project/dao/mock"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/yitter/idgenerator-go/idgen"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupIDGenerator() {
	var options = idgen.NewIdGeneratorOptions(1) // 1 is a worker ID for testing
	idgen.SetIdGenerator(options)
}

func TestCreateGameDetail_CreateSuccess(t *testing.T) {
	setupIDGenerator()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)

	GameDao = mockGameDAO

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
	assert.NotEqual(t, int64(0), resp.GameID, "Expected a new GameID to be generated")
}

func TestCreateGameDetail_UpdateSuccess(t *testing.T) {
	setupIDGenerator()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	mockGameDAO.EXPECT().CreateGameVersionAndUpdateGame(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	req := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 12345,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "My Game V2",
			},
		},
	}

	resp, err := CreateGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
	assert.Equal(t, req.GameDetail.GameID, resp.GameID, "Expected GameID to be the same as in the request")
}

func TestCreateGameDetail_UpdateGameNotFound(t *testing.T) {
	setupIDGenerator()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	mockGameDAO.EXPECT().CreateGameVersionAndUpdateGame(gomock.Any(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound).Times(1)

	req := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 99999,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "My Game V2",
			},
		},
	}

	resp, err := CreateGameDetail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10001", resp.BaseResp.Code, "Expected game not found error code")
	assert.Equal(t, "Game not found", resp.BaseResp.Msg)
}

func TestCreateGameDetail_DaoErrorOnCreate(t *testing.T) {
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
	assert.Equal(t, "500", resp.BaseResp.Code, "Expected internal server error code")
}
