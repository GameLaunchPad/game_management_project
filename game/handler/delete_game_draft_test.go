// file: game/handler/delete_game_draft_test.go
package handler

import (
	"context"
	"testing"

	"github.com/GameLaunchPad/game_management_project/dao"
	"github.com/GameLaunchPad/game_management_project/dao/mock"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestDeleteGameDraft_Success tests the successful deletion of a game draft
func TestDeleteGameDraft_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(101)

	// define expectation: DAO's DeleteGameDraft method is called with correct parameters and returns success
	mockGameDAO.EXPECT().
		DeleteGameDraft(gomock.Any(), gameID).
		Return(nil).
		Times(1)

	req := &game.DeleteGameDraftRequest{
		GameID: int64(gameID),
	}

	resp, err := DeleteGameDraft(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestDeleteGameDraft_GameNotFound tests the scenario where the game to be deleted is not found
func TestDeleteGameDraft_GameNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(999)

	// define expectation: DAO returns gorm.ErrRecordNotFound
	mockGameDAO.EXPECT().
		DeleteGameDraft(gomock.Any(), gameID).
		Return(gorm.ErrRecordNotFound).
		Times(1)

	req := &game.DeleteGameDraftRequest{
		GameID: int64(gameID),
	}

	resp, err := DeleteGameDraft(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10001", resp.BaseResp.Code)
	assert.Equal(t, "Game not found", resp.BaseResp.Msg)
}

// TestDeleteGameDraft_NotDraft tests the scenario where the game version is not a draft
func TestDeleteGameDraft_NotDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(102)

	// define expectation: DAO returns dao.ErrVersionIsNotDraft
	mockGameDAO.EXPECT().
		DeleteGameDraft(gomock.Any(), gameID).
		Return(dao.ErrVersionIsNotDraft).
		Times(1)

	req := &game.DeleteGameDraftRequest{
		GameID: int64(gameID),
	}

	resp, err := DeleteGameDraft(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10003", resp.BaseResp.Code)
	assert.Equal(t, dao.ErrVersionIsNotDraft.Error(), resp.BaseResp.Msg)
}
