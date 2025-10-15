package handler

import (
	"context"
	"testing"

	"github.com/GameLaunchPad/game_management_project/dao/mock"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestUpdateGameDraft_Success tests the successful update of a game draft
func TestUpdateGameDraft_Success(t *testing.T) {
	setupIDGenerator()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	mockGameDAO.EXPECT().UpdateGameDraft(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	req := &game.UpdateGameDraftRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 12345, // GameID不为0
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "My Game V2",
			},
		},
	}

	resp, err := UpdateGameDraft(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestUpdateGameDraft_GameNotFound tests the scenario where the game to be updated does not exist
func TestUpdateGameDraft_GameNotFound(t *testing.T) {
	setupIDGenerator()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	mockGameDAO.EXPECT().UpdateGameDraft(gomock.Any(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound).Times(1)

	req := &game.UpdateGameDraftRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: 99999,
			CpID:   1001,
			GameVersion: &game.GameVersion{
				GameName: "My Game V2",
			},
		},
	}

	resp, err := UpdateGameDraft(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10001", resp.BaseResp.Code)
	assert.Equal(t, "Game not found", resp.BaseResp.Msg)
}

// TestUpdateGameDraft_FailWithZeroGameID tests the failure case when GameID is zero
func TestUpdateGameDraft_FailWithZeroGameID(t *testing.T) {
	req := &game.UpdateGameDraftRequest{
		GameDetail: &game.GameDetailWrite{
			GameID:      0,
			CpID:        1001,
			GameVersion: &game.GameVersion{},
		},
	}

	resp, err := UpdateGameDraft(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "GameID is required")
}
