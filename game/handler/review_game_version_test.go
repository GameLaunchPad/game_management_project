// file: game/handler/review_game_version_test.go
package handler

import (
	"context"
	"testing"

	"github.com/GameLaunchPad/game_management_project/game/dao/mock"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestReviewGameVersion_PassSuccess test the successful scenario of passing a review
func TestReviewGameVersion_PassSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(101)
	versionID := uint64(201)
	expectedStatus := int(game.GameStatus_Published)

	// define expectation: DAO's ReviewGameVersion method is called with correct parameters and returns success
	mockGameDAO.EXPECT().
		ReviewGameVersion(gomock.Any(), gameID, versionID, expectedStatus, "").
		Return(nil).
		Times(1)

	req := &game.ReviewGameVersionRequest{
		GameID:        int64(gameID),
		GameVersionID: int64(versionID),
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestReviewGameVersion_RejectSuccess test the successful scenario of rejecting a review
func TestReviewGameVersion_RejectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(102)
	versionID := uint64(202)
	expectedStatus := int(game.GameStatus_Rejected)

	// define expectation: DAO's ReviewGameVersion method is called with correct parameters and returns success
	mockGameDAO.EXPECT().
		ReviewGameVersion(gomock.Any(), gameID, versionID, expectedStatus, "").
		Return(nil).
		Times(1)

	req := &game.ReviewGameVersionRequest{
		GameID:        int64(gameID),
		GameVersionID: int64(versionID),
		ReviewResult_: game.ReviewResult__Reject,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "200", resp.BaseResp.Code)
}

// TestReviewGameVersion_NotFound tests the scenario where the game or version is not found
func TestReviewGameVersion_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	gameID := uint64(999)
	versionID := uint64(9999)

	// define expectation: DAO's ReviewGameVersion method is called and returns record not found error
	mockGameDAO.EXPECT().
		ReviewGameVersion(gomock.Any(), gameID, versionID, gomock.Any(), gomock.Any()).
		Return(gorm.ErrRecordNotFound).
		Times(1)

	req := &game.ReviewGameVersionRequest{
		GameID:        int64(gameID),
		GameVersionID: int64(versionID),
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10002", resp.BaseResp.Code)
	assert.Equal(t, "Game or Version not found", resp.BaseResp.Msg)
}
