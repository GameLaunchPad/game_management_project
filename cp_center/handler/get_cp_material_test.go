package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	mock_repo "github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetCPMaterial(t *testing.T) {
	ctx := context.Background()

	// 测试场景 1: 成功获取 CP 素材
	t.Run("Success - Get CP Material", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.GetCPMaterialRequest{CpID: 123}
		expectedMaterial := &ddl.GpCpMaterial{
			Id:              1,
			CpId:            123,
			CpName:          "Test CP",
			BusinessLicense: "License ABC",
			Status:          int(cp_center.MaterialStatus_Online),
			CreateTs:        time.Now(),
			ModifyTs:        time.Now(),
		}

		// --- 设定 mock 预期 ---
		// 期望 GetMaterialByCPID 方法被调用 1 次，参数是 123，并成功返回我们预设的 material 对象
		mockRepo.EXPECT().GetMaterialByCPID(ctx, int32(req.CpID)).Return(expectedMaterial, nil).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.GetCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.CPMaterial)
		assert.Equal(t, expectedMaterial.CpName, resp.CPMaterial.CpName)
		assert.Equal(t, int64(expectedMaterial.Id), resp.CPMaterial.MaterialID)
	})

	// 测试场景 2: CP 素材不存在 (gorm.ErrRecordNotFound)
	t.Run("Failure - Material Not Found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.GetCPMaterialRequest{CpID: 404}

		// --- 设定 mock 预期 ---
		// 期望 GetMaterialByCPID 方法被调用 1 次，但这次返回 gorm.ErrRecordNotFound 错误
		mockRepo.EXPECT().GetMaterialByCPID(ctx, int32(req.CpID)).Return(nil, gorm.ErrRecordNotFound).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.GetCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "cp material not found")
	})

	// 测试场景 3: 数据库发生其他错误
	t.Run("Failure - Other Database Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.GetCPMaterialRequest{CpID: 500}
		dbError := errors.New("database connection lost")

		// --- 设定 mock 预期 ---
		// 期望 GetMaterialByCPID 方法被调用 1 次，返回一个通用的数据库错误
		mockRepo.EXPECT().GetMaterialByCPID(ctx, int32(req.CpID)).Return(nil, dbError).Times(1)

		// --- 执行被测函数 ---
		resp, err := h.GetCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbError) // 确认返回的错误就是我们模拟的那个
	})

	// 测试场景 4: 输入参数无效 (CpID <= 0)
	t.Run("Failure - Invalid Parameter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mock_repo.NewMockICPMaterialRepo(ctrl)
		h := handler.NewCPMaterialHandler(mockRepo)

		req := &cp_center.GetCPMaterialRequest{CpID: 0}

		// --- 设定 mock 预期 ---
		// 因为参数校验在数据库调用之前，所以我们期望 Repo 的任何方法都**不**被调用。
		// gomock 在这里不需要写 mockRepo.EXPECT()。

		// --- 执行被测函数 ---
		resp, err := h.GetCPMaterial(ctx, req)

		// --- 断言结果 ---
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid parameter: cp_id must be positive")
	})
}
