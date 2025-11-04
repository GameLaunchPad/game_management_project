// file: game/handler/review_game_version_test.go
package handler

import (
	"context"
	"errors"
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

// TestReviewGameVersion_InvalidGameID tests the failure case when GameID is 0 or negative
func TestReviewGameVersion_InvalidGameID(t *testing.T) {
	req := &game.ReviewGameVersionRequest{
		GameID:        0,
		GameVersionID: 201,
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "Invalid GameID or GameVersionID")
}

// TestReviewGameVersion_NegativeGameID tests the failure case when GameID is negative
func TestReviewGameVersion_NegativeGameID(t *testing.T) {
	req := &game.ReviewGameVersionRequest{
		GameID:        -1,
		GameVersionID: 201,
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "Invalid GameID or GameVersionID")
}

// TestReviewGameVersion_InvalidVersionID tests the failure case when GameVersionID is 0 or negative
func TestReviewGameVersion_InvalidVersionID(t *testing.T) {
	req := &game.ReviewGameVersionRequest{
		GameID:        101,
		GameVersionID: 0,
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "Invalid GameID or GameVersionID")
}

// TestReviewGameVersion_NegativeVersionID tests the failure case when GameVersionID is negative
func TestReviewGameVersion_NegativeVersionID(t *testing.T) {
	req := &game.ReviewGameVersionRequest{
		GameID:        101,
		GameVersionID: -1,
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "400", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "Invalid GameID or GameVersionID")
}

// TestReviewGameVersion_OtherError tests the failure case when DAO returns other errors (not RecordNotFound)
func TestReviewGameVersion_OtherError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGameDAO := mock.NewMockIGameDAO(ctrl)
	GameDao = mockGameDAO

	otherError := errors.New("database connection error")
	mockGameDAO.EXPECT().
		ReviewGameVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(otherError).
		Times(1)

	req := &game.ReviewGameVersionRequest{
		GameID:        101,
		GameVersionID: 201,
		ReviewResult_: game.ReviewResult__Pass,
	}

	resp, err := ReviewGameVersion(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "500", resp.BaseResp.Code)
	assert.Contains(t, resp.BaseResp.Msg, "Failed to update game version status")
}
