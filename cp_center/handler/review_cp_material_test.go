package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// mustReviewRemark 是一个辅助函数，用于创建 ReviewRemark 指针
func mustReviewRemark(remark string, operator string) *cp_center.ReviewRemark {
	return &cp_center.ReviewRemark{
		Remark:   remark,
		Operator: operator,
	}
}

// updateMapMatcher 是一个 gomock 匹配器，用于验证更新 map 的内容
type updateMapMatcher struct {
	expectedStatus   int
	expectedRemark   string
	expectedOperator string
}

// NewUpdateMapMatcher 创建一个新的匹配器
func NewUpdateMapMatcher(status int, remark, operator string) gomock.Matcher {
	return &updateMapMatcher{
		expectedStatus:   status,
		expectedRemark:   remark,
		expectedOperator: operator,
	}
}

// Matches 检查传入的 map 是否符合预期
func (m *updateMapMatcher) Matches(x interface{}) bool {
	updates, ok := x.(map[string]interface{})
	if !ok {
		return false
	}

	status, ok := updates["status"].(int)
	if !ok || status != m.expectedStatus {
		return false
	}

	remark, ok := updates["review_comment"].(string)
	if !ok || remark != m.expectedRemark {
		return false
	}

	operator, ok := updates["operator"].(string)
	if !ok || operator != m.expectedOperator {
		return false
	}

	_, ok = updates["modify_ts"].(time.Time)
	return ok // 确保 modify_ts 字段存在且类型正确
}

// String 描述了匹配器
func (m *updateMapMatcher) String() string {
	return "is an update map with correct status, remark, and operator"
}

func TestCPMaterialHandler_ReviewCPMaterial(t *testing.T) {
	type testHandler struct {
		*handler.CPMaterialHandler
		MockMaterialRepo *mocks.MockICPMaterialRepo
		MockCPRepo       *mocks.MockICPRepo // 添加 MockCPRepo
	}

	// 辅助函数来构建 handler 和 mock
	setup := func(t *testing.T) (*testHandler, *gomock.Controller) {
		ctrl := gomock.NewController(t)
		mockMaterialRepo := mocks.NewMockICPMaterialRepo(ctrl)
		mockCPRepo := mocks.NewMockICPRepo(ctrl) // 创建 MockCPRepo

		// 传递两个 mock 依赖
		realHandler := handler.NewCPMaterialHandler(mockMaterialRepo, mockCPRepo)

		return &testHandler{
			CPMaterialHandler: realHandler,
			MockMaterialRepo:  mockMaterialRepo,
			MockCPRepo:        mockCPRepo, // 存储 MockCPRepo
		}, ctrl
	}

	// 预期成功响应
	successResp := &cp_center.ReviewCPMaterialResponse{
		BaseResp: &common.BaseResp{
			Code: "0",
			Msg:  "success",
		},
	}

	// 定义测试用例
	tests := []struct {
		name      string
		req       *cp_center.ReviewCPMaterialRequest
		mockSetup func(mockRepo *mocks.MockICPMaterialRepo) // mockSetup 保持不变，因为 ReviewCPMaterial 只用到了 MaterialRepo
		want      *cp_center.ReviewCPMaterialResponse
		wantErr   assert.ErrorAssertionFunc
		errMsg    string
	}{
		{
			name: "Error: Invalid MaterialID",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    0, // 无效ID
				ReviewResult_: cp_center.ReviewResult__Pass,
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 期望没有数据库调用
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "invalid parameter: material_id is required",
		},
		{
			name: "Error: Unset ReviewResult",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    1,
				ReviewResult_: cp_center.ReviewResult__Unset, // 无效结果
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 期望没有数据库调用
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "invalid parameter: review_result must be Pass or Reject",
		},
		{
			name: "Error: Material Not Found",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    404,
				ReviewResult_: cp_center.ReviewResult__Pass,
				ReviewRemark:  mustReviewRemark("test", "admin"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(404)).
					Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "material not found",
		},
		{
			name: "Error: GetMaterialByID DB Error",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    500,
				ReviewResult_: cp_center.ReviewResult__Pass,
				ReviewRemark:  mustReviewRemark("test", "admin"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(500)).
					Return(nil, errors.New("db connection error"))
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "db connection error",
		},
		{
			name: "Success: Review Pass",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    1,
				ReviewResult_: cp_center.ReviewResult__Pass,
				ReviewRemark:  mustReviewRemark("Looks good", "admin-pass"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 1. Get
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(1)).
					Return(&ddl.GpCpMaterial{Id: 1, CpId: 10}, nil)

				// 2. Update
				matcher := NewUpdateMapMatcher(3, "Looks good", "admin-pass")
				mockRepo.EXPECT().
					UpdateMaterial(gomock.Any(), int64(1), matcher).
					Return(int64(1), nil) // 1 row affected
			},
			want:    successResp,
			wantErr: assert.NoError,
		},
		{
			name: "Success: Review Reject",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    2,
				ReviewResult_: cp_center.ReviewResult__Reject,
				ReviewRemark:  mustReviewRemark("Missing info", "admin-reject"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 1. Get
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(2)).
					Return(&ddl.GpCpMaterial{Id: 2, CpId: 11}, nil)

				// 2. Update
				matcher := NewUpdateMapMatcher(4, "Missing info", "admin-reject")
				mockRepo.EXPECT().
					UpdateMaterial(gomock.Any(), int64(2), matcher).
					Return(int64(1), nil) // 1 row affected
			},
			want:    successResp,
			wantErr: assert.NoError,
		},
		{
			name: "Error: UpdateMaterial DB Error",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    3,
				ReviewResult_: cp_center.ReviewResult__Pass,
				ReviewRemark:  mustReviewRemark("test", "admin"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 1. Get
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(3)).
					Return(&ddl.GpCpMaterial{Id: 3}, nil)

				// 2. Update (fails)
				mockRepo.EXPECT().
					UpdateMaterial(gomock.Any(), int64(3), gomock.Any()).
					Return(int64(0), errors.New("update failed error"))
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "update failed error",
		},
		{
			name: "Error: UpdateMaterial Zero Rows Affected",
			req: &cp_center.ReviewCPMaterialRequest{
				MaterialID:    4,
				ReviewResult_: cp_center.ReviewResult__Pass,
				ReviewRemark:  mustReviewRemark("test", "admin"),
			},
			mockSetup: func(mockRepo *mocks.MockICPMaterialRepo) {
				// 1. Get
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(4)).
					Return(&ddl.GpCpMaterial{Id: 4}, nil)

				// 2. Update (returns 0 rows)
				mockRepo.EXPECT().
					UpdateMaterial(gomock.Any(), int64(4), gomock.Any()).
					Return(int64(0), nil) // 0 rows affected
			},
			want:    nil,
			wantErr: assert.Error,
			errMsg:  "update failed, zero rows affected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ctrl := setup(t)
			defer ctrl.Finish()

			// 设置 mock 预期
			if tt.mockSetup != nil {
				tt.mockSetup(h.MockMaterialRepo)
			}

			// 执行被测函数
			got, err := h.CPMaterialHandler.ReviewCPMaterial(context.Background(), tt.req)
			// 验证错误
			tt.wantErr(t, err)
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error(), "Error message mismatch")
			}

			// 验证响应
			assert.Equal(t, tt.want, got, "Response mismatch")
		})
	}
}
