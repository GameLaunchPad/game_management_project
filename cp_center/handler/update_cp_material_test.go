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

func TestUpdateCPMaterial(t *testing.T) {
	ctx := context.Background()
	// 准备一个可编辑状态的原始素材
	editableMaterial := &ddl.GpCpMaterial{Id: 1, CpId: 123, Status: 1} // 1-草稿
	// 准备一个不可编辑状态的原始素材
	uneditableMaterial := &ddl.GpCpMaterial{Id: 2, CpId: 124, Status: 2} // 2-审核中

	// 准备一个基础的、合法的请求对象
	baseReq := &cp_center.UpdateCPMaterialRequest{
		MaterialID: 1,
		CpMaterial: &cp_center.CPMaterial{
			CpName:           "Updated CP Name",
			BusinessLicenses: "Updated License",
		},
		SubmitMode: cp_center.SubmitMode_SubmitDraft,
	}

	// 测试场景 1: 成功更新为草稿
	t.Run("Success - Update as Draft", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, baseReq.MaterialID).Return(editableMaterial, nil).Times(1)
		mockRepo.EXPECT().UpdateMaterial(ctx, baseReq.MaterialID, gomock.Any()).Return(int64(1), nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.UpdateCPMaterial(ctx, baseReq)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.BaseResp.Code)
	})

	// 测试场景 2: 成功更新并提交审核
	t.Run("Success - Submit for Review", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := *baseReq // 复制基础请求
		req.SubmitMode = cp_center.SubmitMode_SubmitReview

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(editableMaterial, nil).Times(1)
		mockRepo.EXPECT().UpdateMaterial(ctx, req.MaterialID, gomock.Any()).Return(int64(1), nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.UpdateCPMaterial(ctx, &req)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	// 测试场景 3: 素材不存在
	t.Run("Failure - Material Not Found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := *baseReq
		req.MaterialID = 999

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.UpdateCPMaterial(ctx, &req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "material not found")
	})

	// 测试场景 4: 业务逻辑错误 - 更新一个正在审核或已发布的素材
	t.Run("Failure - Business Logic Error - Cannot Update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := *baseReq
		req.MaterialID = 2

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, req.MaterialID).Return(uneditableMaterial, nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.UpdateCPMaterial(ctx, &req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "cannot update material that is in review or online")
	})

	// 测试场景 5: 数据库更新操作失败
	t.Run("Failure - Database Update Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		dbError := errors.New("db connection failed")

		// --- 设定 mock 预期 ---
		mockRepo.EXPECT().GetMaterialByID(ctx, baseReq.MaterialID).Return(editableMaterial, nil).Times(1)
		mockRepo.EXPECT().UpdateMaterial(ctx, baseReq.MaterialID, gomock.Any()).Return(int64(0), dbError).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.UpdateCPMaterial(ctx, baseReq)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbError)
	})

	// 测试场景 6: 无效参数
	t.Run("Failure - Invalid Parameters", func(t *testing.T) {
		testCases := map[string]*cp_center.UpdateCPMaterialRequest{
			"MaterialID is zero":  {MaterialID: 0, CpMaterial: baseReq.CpMaterial, SubmitMode: baseReq.SubmitMode},
			"CpMaterial is nil":   {MaterialID: 1, CpMaterial: nil, SubmitMode: baseReq.SubmitMode},
			"SubmitMode is unset": {MaterialID: 1, CpMaterial: baseReq.CpMaterial, SubmitMode: cp_center.SubmitMode_Unset},
			"CpName is empty":     {MaterialID: 1, CpMaterial: &cp_center.CPMaterial{BusinessLicenses: "abc"}, SubmitMode: baseReq.SubmitMode},
		}

		for name, req := range testCases {
			t.Run(name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
				h := handler.NewCPMaterialHandler(mockRepo)

				// 参数校验失败，不应有数据库调用
				resp, err := h.UpdateCPMaterial(ctx, req)

				assert.Nil(t, resp)
				assert.Error(t, err)
			})
		}
	})
}
