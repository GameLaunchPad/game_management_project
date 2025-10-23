package handler_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/handler"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/common"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestCPMaterialHandler_GetCPMaterial(t *testing.T) {
	// 1. 初始化 Gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. 创建 Mock Repository 实例
	mockRepo := mocks.NewMockICPMaterialRepo(ctrl)

	// 3. 创建被测试的 Handler 实例，并注入 Mock
	h := &handler.CPMaterialHandler{
		MaterialRepo: mockRepo,
	}

	// 4. 定义通用的上下文和模拟数据
	ctx := context.Background()
	mockTime := time.Now()
	mockMaterial := &ddl.GpCpMaterial{
		Id:              101,
		CpId:            1,
		CpIcon:          "test_icon.png",
		CpName:          "Test CP",
		BusinessLicense: "license1.jpg",
		Website:         "https://test.com",
		Status:          1, // 假设 1 = 审核通过
		ReviewComment:   "ok",
		CreateTs:        mockTime,
		ModifyTs:        mockTime,
	}

	// 5. 定义测试用例
	type args struct {
		ctx context.Context
		req *cp_center.GetCPMaterialRequest
	}
	type want struct {
		resp *cp_center.GetCPMaterialResponse
		err  error
	}
	type mockSetup func()

	tests := []struct {
		name      string
		args      args
		want      want
		mockSetup mockSetup
	}{
		{
			name: "Success - Get material successfully",
			args: args{
				ctx: ctx,
				req: &cp_center.GetCPMaterialRequest{
					CpID:       1,
					MaterialID: 101,
				},
			},
			want: want{
				resp: &cp_center.GetCPMaterialResponse{
					BaseResp: &common.BaseResp{Code: "0", Msg: "success"},
					CPMaterial: &cp_center.CPMaterial{
						MaterialID:         101,
						CpID:               1,
						CpIcon:             "test_icon.png",
						CpName:             "Test CP",
						VerificationImages: nil, // 根据handler逻辑，这里是nil
						BusinessLicenses:   "license1.jpg",
						Website:            "https://test.com",
						Status:             cp_center.MaterialStatus(1), // 对应 mockMaterial.Status
						ReviewComment:      "ok",
						CreateTime:         mockTime.Unix(),
						ModifyTime:         mockTime.Unix(),
					},
				},
				err: nil,
			},
			mockSetup: func() {
				// 定义 Mock 期望：当 GetMaterialByID 被调用时，返回模拟数据
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(101)).
					Return(mockMaterial, nil).
					Times(1) // 期望被调用1次
			},
		},
		{
			name: "Failure - Invalid parameter cp_id",
			args: args{
				ctx: ctx,
				req: &cp_center.GetCPMaterialRequest{
					CpID:       0, // 非法的 CpID
					MaterialID: 101,
				},
			},
			want: want{
				resp: &cp_center.GetCPMaterialResponse{
					BaseResp: &common.BaseResp{Code: "500", Msg: "invalid parameter: cp_id must be positive: "},
				},
				err: nil, // Handler 内部处理了错误，返回给 kitex 的 error 为 nil
			},
			mockSetup: func() {
				// 参数校验失败，不应该调用 Repo
				mockRepo.EXPECT().GetMaterialByID(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "Failure - Repository returns error",
			args: args{
				ctx: ctx,
				req: &cp_center.GetCPMaterialRequest{
					CpID:       1,
					MaterialID: 102, // 假设这个ID会导致错误
				},
			},
			want: want{
				resp: &cp_center.GetCPMaterialResponse{
					BaseResp: &common.BaseResp{Code: "500", Msg: "record not found"},
				},
				err: nil,
			},
			mockSetup: func() {
				// 定义 Mock 期望：当 GetMaterialByID 被调用时，返回一个错误
				mockRepo.EXPECT().
					GetMaterialByID(gomock.Any(), int64(102)).
					Return(nil, errors.New("record not found")).
					Times(1)
			},
		},
	}

	// 6. 循环执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// (a) 执行每个测试用例的 Mock 设置
			tt.mockSetup()

			// (b) 调用被测试的方法
			got, err := h.GetCPMaterial(tt.args.ctx, tt.args.req)

			// (c) 断言返回的 error
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("GetCPMaterial() error = %v, wantErr %v", err, tt.want.err)
				return
			}

			// (d) 断言返回的 response
			// 注意：BaseResp 和 CPMaterial 都是指针，DeepEqual 会比较它们指向的内容
			if !reflect.DeepEqual(got, tt.want.resp) {
				t.Errorf("GetCPMaterial() got = %v, want %v", got, tt.want.resp)
			}
		})
	}
}
