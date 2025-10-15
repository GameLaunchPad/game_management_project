package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	mock_repo "github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestReviewCPMaterial(t *testing.T) {
	ctx := context.Background()
	existingMaterial := &ddl.GpCpMaterial{Id: 1, CpId: 123, Status: 2} // 假设原始状态是审核中

	// 测试场景 1: 成功审核通过
	t.Run("Success - Review Pass", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.ReviewCPMaterialRequest{
			MaterialID:    1,
			ReviewResult_: cp_center.ReviewResult__Pass,
		}

		// --- 设定 mock 预期 ---
		// 1. 期望 GetMaterialByID 被调用，且成功返回一个存在的素材
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(existingMaterial, nil).Times(1)
		// 2. 期望 UpdateMaterial 被调用，且成功更新了 1 行
		mockRepo.EXPECT().UpdateMaterial(ctx, req.MaterialID, gomock.Any()).Return(int64(1), nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.ReviewCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.BaseResp.Code)
	})

	// 测试场景 2: 成功审核拒绝
	t.Run("Success - Review Reject", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.ReviewCPMaterialRequest{
			MaterialID:    1,
			ReviewResult_: cp_center.ReviewResult__Reject,
		}

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(existingMaterial, nil).Times(1)
		mockRepo.EXPECT().UpdateMaterial(ctx, req.MaterialID, gomock.Any()).Return(int64(1), nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.ReviewCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	// 测试场景 3: 素材不存在
	t.Run("Failure - Material Not Found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.ReviewCPMaterialRequest{MaterialID: 999, ReviewResult_: cp_center.ReviewResult__Pass}

		// --- 设定 mock 预期 ---
		// 期望 GetMaterialByID 返回 gorm.ErrRecordNotFound
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(nil, gorm.ErrRecordNotFound).Times(1)
		// 不期望 UpdateMaterial 被调用

		// --- 执行被测函数 ---
		resp, err := h.ReviewCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "material not found")
	})

	// 测试场景 4: 更新数据库失败
	t.Run("Failure - Update Database Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.ReviewCPMaterialRequest{MaterialID: 1, ReviewResult_: cp_center.ReviewResult__Pass}
		dbError := errors.New("update failed")

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(existingMaterial, nil).Times(1)
		// 期望 UpdateMaterial 调用失败
		mockRepo.EXPECT().UpdateMaterial(ctx, req.MaterialID, gomock.Any()).Return(int64(0), dbError).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.ReviewCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbError)
	})

	// 测试场景 5: 更新影响 0 行
	t.Run("Failure - Update Affects Zero Rows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.ReviewCPMaterialRequest{MaterialID: 1, ReviewResult_: cp_center.ReviewResult__Pass}

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(existingMaterial, nil).Times(1)
		// 期望 UpdateMaterial 返回成功，但影响行数为 0
		mockRepo.EXPECT().UpdateMaterial(ctx, req.MaterialID, gomock.Any()).Return(int64(0), nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.ReviewCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "update failed, zero rows affected")
	})

	// 测试场景 6: 无效参数
	t.Run("Failure - Invalid Parameters", func(t *testing.T) {
		// 子测试用例
		testCases := map[string]*cp_center.ReviewCPMaterialRequest{
			"MaterialID is zero":    {MaterialID: 0, ReviewResult_: cp_center.ReviewResult__Pass},
			"ReviewResult is unset": {MaterialID: 1, ReviewResult_: cp_center.ReviewResult__Unset},
		}

		for name, req := range testCases {
			t.Run(name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
				h := handler.NewCPMaterialHandler(mockRepo)

				// --- 设定 mock 预期 ---
				// 参数校验失败，不应该有任何数据库调用

				// --- 执行被测函数 ---
				resp, err := h.ReviewCPMaterial(ctx, req)

				// --- 断言结果 ---
				assert.Nil(t, resp)
				assert.Error(t, err)
			})
		}
	})
}
