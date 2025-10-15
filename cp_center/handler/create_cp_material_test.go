// Package handler_test 是 handler 包的测试包
package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	mock_repo "github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestCreateCPMaterial 使用 gomock 对 CreateCPMaterial 方法进行单元测试
func TestCreateCPMaterial(t *testing.T) {
	// 初始化测试上下文
	ctx := context.Background()

	// 准备一个基础的、合法的请求对象，方便在各个测试用例中复用
	baseReq := &cp_center.CreateCPMaterialRequest{
		CPMaterial: &cp_center.CPMaterial{
			CpID:             123,
			CpName:           "Test CP",
			BusinessLicenses: "License ABC",
		},
		SubmitMode: cp_center.SubmitMode_SubmitDraft,
	}

	// 测试场景 1: 成功创建草稿
	t.Run("Success - Create as Draft", func(t *testing.T) {
		// gomock.Controller 是 mock 的核心，管理 mock 对象的生命周期和期望
		ctrl := gomock.NewController(t)
		// 创建一个 mockRepo 实例 (这就是“假的”Repo)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)

		// 将“假的”Repo 注入到我们的 handler 中
		h := handler.NewCPMaterialHandler(mockRepo)

		// --- 设定 mock 预期 ---
		// 我们期望 mockRepo 的 CreateMaterial 方法会被精确调用 1 次。
		// gomock.Any() 表示我们不关心传入的具体参数是什么，只要类型正确即可。
		// 当它被调用时，我们让它假装成功，返回 nil。
		mockRepo.EXPECT().CreateMaterial(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.CreateCPMaterial(ctx, baseReq)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.BaseResp.Code)
		assert.Equal(t, "创建成功", resp.BaseResp.Msg)
		assert.Equal(t, baseReq.CPMaterial.CpID, resp.CpID)
	})

	// 测试场景 2: 数据库创建失败
	t.Run("Failure - Database Create Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		// 模拟一个数据库错误
		dbError := errors.New("unique constraint failed")

		// --- 设定 mock 预期 ---
		// 我们期望 CreateMaterial 方法被调用 1 次，但这次让它假装失败，返回我们模拟的 dbError。
		mockRepo.EXPECT().CreateMaterial(gomock.Any(), gomock.Any()).Return(dbError).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.CreateCPMaterial(ctx, baseReq)

		// --- 断言结果 ---
		assert.Nil(t, resp)                                                // 失败时，响应体应为 nil
		assert.Error(t, err)                                               // 应该返回错误
		assert.ErrorContains(t, err, "failed to create cp material in db") // 检查错误信息是否被正确包装
		assert.ErrorIs(t, err, dbError)                                    // 检查错误链中是否包含了原始的数据库错误
	})

	// 测试场景 3: 输入参数无效
	t.Run("Failure - Invalid Arguments", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		// 定义一系列不合法的请求来进行测试
		testCases := map[string]*cp_center.CreateCPMaterialRequest{
			"CPMaterial is nil": {CPMaterial: nil},
			"CpID is zero":      {CPMaterial: &cp_center.CPMaterial{CpName: "Test"}},
			"CpName is empty":   {CPMaterial: &cp_center.CPMaterial{CpID: 123}},
		}

		for name, req := range testCases {
			t.Run(name, func(t *testing.T) {
				// --- 设定 mock 预期 ---
				// 在参数校验失败的情况下，代码逻辑不应该调用数据库。
				// 所以，我们不设定任何 mockRepo.EXPECT()。
				// 如果 handler 的代码逻辑有误，意外调用了 CreateMaterial，gomock 会因为没有匹配的预期而让测试失败。

				// --- 执行被测函数 ---
				resp, err := h.CreateCPMaterial(ctx, req)

				// --- 断言结果 ---
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.ErrorContains(t, err, "invalid arguments")
			})
		}
	})
}
