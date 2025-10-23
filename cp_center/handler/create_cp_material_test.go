package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/GameLaunchPad/game_management_project/cp_center/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/cp_center/kitex_gen/cp_center"
	"github.com/GameLaunchPad/game_management_project/cp_center/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// setupIDGenerator 初始化 idgen，以便 idgen.NextId() 可以工作
// 在实际测试中，idgen 通常在 main 包或 init() 中初始化。
func setupIDGenerator() {
	// 设置一个简单的 WorkerId，仅用于测试目的
	options := idgen.NewIdGeneratorOptions(1)
	idgen.SetIdGenerator(options)
}

// materialMatcher 是一个自定义的 gomock 匹配器
// 它用于匹配 GpCpMaterial 对象，但忽略动态生成的字段（如 Id, CreateTs, ModifyTs）
type materialMatcher struct {
	expected *ddl.GpCpMaterial
}

// Matches 实现了 gomock.Matcher 接口
func (m materialMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*ddl.GpCpMaterial)
	if !ok {
		return false
	}

	// 比较我们关心的核心业务字段
	return actual.CpId == m.expected.CpId &&
		actual.CpName == m.expected.CpName &&
		actual.CpIcon == m.expected.CpIcon &&
		actual.BusinessLicense == m.expected.BusinessLicense &&
		actual.Website == m.expected.Website &&
		actual.VerificationImages == m.expected.VerificationImages &&
		actual.Status == m.expected.Status
}

// String 实现了 gomock.Matcher 接口
func (m materialMatcher) String() string {
	return fmt.Sprintf("is a GpCpMaterial matching CpId=%d, CpName=%s, Status=%d",
		m.expected.CpId, m.expected.CpName, m.expected.Status)
}

// MaterialEq 返回一个 materialMatcher 实例
func MaterialEq(expected *ddl.GpCpMaterial) gomock.Matcher {
	return materialMatcher{expected: expected}
}

// cpMatcher 是一个自定义的 gomock 匹配器
// 它用于匹配 GpCp 对象，忽略动态时间戳
type cpMatcher struct {
	expected *ddl.GpCp
}

// Matches 实现了 gomock.Matcher 接口
func (m cpMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*ddl.GpCp)
	if !ok {
		return false
	}

	// 比较核心字段
	// 我们不能比较 NewestMaterialId，因为它是在 CreateMaterial 之后才确定的
	// 但在 CreateCP 场景下，NewestMaterialId 来自 material.Id，是动态的
	// 所以我们只比较 CpId 和 CpName
	// 我们可以让 MaterialEq 匹配器在匹配时捕获 ID，但这会使测试变得复杂。
	// 让我们假设 Id 是从 material 正确传递过来的。
	return actual.Id == m.expected.Id &&
		actual.CpName == m.expected.CpName &&
		actual.VerifyStatus == m.expected.VerifyStatus
}

// String 实现了 gomock.Matcher 接口
func (m cpMatcher) String() string {
	return fmt.Sprintf("is a GpCp matching Id=%d, CpName=%s",
		m.expected.Id, m.expected.CpName)
}

// CPEq 返回一个 cpMatcher 实例
func CPEq(expected *ddl.GpCp) gomock.Matcher {
	return cpMatcher{expected: expected}
}

// TestCreateCPMaterial 包含所有 CreateCPMaterial 的子测试
func TestCreateCPMaterial(t *testing.T) {
	// 确保 idgen 已初始化
	setupIDGenerator()

	// --- 通用设置 ---
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMaterialRepo := mocks.NewMockICPMaterialRepo(mockCtrl)
	mockCPRepo := mocks.NewMockICPRepo(mockCtrl)

	handler := &CPMaterialHandler{
		MaterialRepo: mockMaterialRepo,
		CPRepo:       mockCPRepo,
	}

	// --- 通用测试数据 ---
	baseReq := &cp_center.CreateCPMaterialRequest{
		CPMaterial: &cp_center.CPMaterial{
			CpID:               1001,
			CpIcon:             "http://icon.url/icon.png",
			CpName:             "Test CP",
			VerificationImages: []string{"http://img.url/img1.png"},
			BusinessLicenses:   "http://img.url/license.png",
			Website:            "http://test-cp.com",
		},
		SubmitMode: cp_center.SubmitMode_SubmitReview,
	}

	expectedMaterialReview := &ddl.GpCpMaterial{
		CpId:               1001,
		CpIcon:             "http://icon.url/icon.png",
		CpName:             "Test CP",
		VerificationImages: `["http://img.url/img1.png"]`, // JSON 序列化后的
		BusinessLicense:    "http://img.url/license.png",
		Website:            "http://test-cp.com",
		Status:             int(cp_center.MaterialStatus_Reviewing),
	}

	expectedCPNew := &ddl.GpCp{
		Id:     1001,
		CpName: "Test CP",
		// NewestMaterialId 在测试中是动态的，无法预先确定
		OnlineMaterialId: 0,
		VerifyStatus:     0,
	}

	// --- 子测试 ---

	t.Run("Success_NewCP_SubmitReview", func(t *testing.T) {
		// 1. Mock 期望
		// 1.1 CreateMaterial 成功
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(expectedMaterialReview)).
			// 我们需要捕获动态生成的 ID，以便在后续步骤中使用
			// 通过 Do 来模拟 Id 的生成并设置它
			Do(func(ctx context.Context, mat *ddl.GpCpMaterial) {
				// 模拟数据库或 idgen 分配了一个 ID
				if mat.Id == 0 {
					mat.Id = 9999 // 分配一个固定的测试 ID
				}
			}).
			Return(nil)

		// 1.2 GetCPByID 返回 "Not Found"，触发创建逻辑
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1001)).
			Return(nil, gorm.ErrRecordNotFound)

		// 1.3 CreateCP 成功
		// 我们期望的 CP 应该包含来自 material 的动态 ID (9999)
		expectedCPNew.NewestMaterialId = 9999 // 设置期望的 ID
		mockCPRepo.EXPECT().
			CreateCP(ctx, CPEq(expectedCPNew)). // 使用自定义匹配器
			Return(nil)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, baseReq)

		// 3. 断言
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.BaseResp.Code)
		assert.Equal(t, "创建成功", resp.BaseResp.Msg)
		assert.Equal(t, int64(1001), resp.CpID)
	})

	t.Run("Success_ExistingCP_SubmitDraft", func(t *testing.T) {
		// 准备草稿模式的请求和期望
		draftReq := &cp_center.CreateCPMaterialRequest{
			CPMaterial: baseReq.CPMaterial, // 复用基础数据
			SubmitMode: cp_center.SubmitMode_SubmitDraft,
		}
		expectedMaterialDraft := *expectedMaterialReview // 复制
		expectedMaterialDraft.Status = int(cp_center.MaterialStatus_Draft)

		// 1. Mock 期望
		// 1.1 CreateMaterial 成功 (草稿模式)
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(&expectedMaterialDraft)).
			Do(func(ctx context.Context, mat *ddl.GpCpMaterial) {
				mat.Id = 8888 // 分配一个不同的 ID
			}).
			Return(nil)

		// 1.2 GetCPByID 成功返回一个已存在的 CP
		existingCP := &ddl.GpCp{
			Id:               1001,
			CpName:           "Old CP Name", // 旧名称
			NewestMaterialId: 7777,          // 旧的资质 ID
			OnlineMaterialId: 7777,
			VerifyStatus:     1,
			CreateTs:         time.Now().Add(-24 * time.Hour),
			ModifyTs:         time.Now().Add(-24 * time.Hour),
		}
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1001)).
			Return(existingCP, nil)

		// 1.3 UpdateCP 成功
		expectedUpdates := map[string]interface{}{
			"newest_material_id": uint64(8888), // 新的资质 ID
			"cp_name":            "Test CP",    // 新的 CP 名称
		}
		mockCPRepo.EXPECT().
			UpdateCP(ctx, int64(1001), expectedUpdates).
			Return(nil)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, draftReq)

		// 3. 断言
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "0", resp.BaseResp.Code)
		assert.Equal(t, int64(1001), resp.CpID)
		assert.Equal(t, int64(8888), resp.MaterialID)
	})

	t.Run("Fail_Validation_NilMaterial", func(t *testing.T) {
		req := &cp_center.CreateCPMaterialRequest{
			CPMaterial: nil, // 故意设置为 nil
		}
		resp, err := handler.CreateCPMaterial(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "CPMaterial cannot be nil")
	})

	t.Run("Fail_Validation_EmptyCpName", func(t *testing.T) {
		req := &cp_center.CreateCPMaterialRequest{
			CPMaterial: &cp_center.CPMaterial{
				CpID:   1001,
				CpName: "   ", // 空白名称
			},
		}
		resp, err := handler.CreateCPMaterial(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "CpName is required")
	})

	t.Run("Fail_CreateMaterial_DBError", func(t *testing.T) {
		dbErr := errors.New("database connection failed")

		// 1. Mock 期望
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(expectedMaterialReview)).
			Return(dbErr)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, baseReq)

		// 3. 断言
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to create cp material in db")
		assert.True(t, errors.Is(err, dbErr)) // 检查错误链
	})

	t.Run("Fail_CheckCP_DBError", func(t *testing.T) {
		dbErr := errors.New("database query failed")

		// 1. Mock 期望
		// 1.1 CreateMaterial 成功
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(expectedMaterialReview)).
			Return(nil)

		// 1.2 GetCPByID 失败 (不是 gorm.ErrRecordNotFound)
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1001)).
			Return(nil, dbErr)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, baseReq)

		// 3. 断言
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to check if cp exists")
		assert.True(t, errors.Is(err, dbErr))
	})

	t.Run("Fail_CreateCP_DBError", func(t *testing.T) {
		dbErr := errors.New("failed to insert new cp")

		// 1. Mock 期望
		// 1.1 CreateMaterial 成功
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(expectedMaterialReview)).
			Do(func(ctx context.Context, mat *ddl.GpCpMaterial) {
				mat.Id = 9999 // 模拟 ID
			}).
			Return(nil)

		// 1.2 GetCPByID 返回 "Not Found"
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1001)).
			Return(nil, gorm.ErrRecordNotFound)

		// 1.3 CreateCP 失败
		expectedCPNew.NewestMaterialId = 9999
		mockCPRepo.EXPECT().
			CreateCP(ctx, CPEq(expectedCPNew)).
			Return(dbErr)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, baseReq)

		// 3. 断言
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to create cp in db")
		assert.True(t, errors.Is(err, dbErr))
	})

	t.Run("Fail_UpdateCP_DBError", func(t *testing.T) {
		dbErr := errors.New("failed to update existing cp")

		// 1. Mock 期望
		// 1.1 CreateMaterial 成功
		mockMaterialRepo.EXPECT().
			CreateMaterial(ctx, MaterialEq(expectedMaterialReview)).
			Do(func(ctx context.Context, mat *ddl.GpCpMaterial) {
				mat.Id = 8888
			}).
			Return(nil)

		// 1.2 GetCPByID 成功返回
		existingCP := &ddl.GpCp{Id: 1001, CpName: "Old Name"}
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1001)).
			Return(existingCP, nil)

		// 1.3 UpdateCP 失败
		expectedUpdates := map[string]interface{}{
			"newest_material_id": uint64(8888),
			"cp_name":            "Test CP",
		}
		mockCPRepo.EXPECT().
			UpdateCP(ctx, int64(1001), expectedUpdates).
			Return(dbErr)

		// 2. 执行
		resp, err := handler.CreateCPMaterial(ctx, baseReq)

		// 3. 断言
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to update existing cp")
		assert.True(t, errors.Is(err, dbErr))
	})
}

// 确保 CPHandler 实现了 gomock 所需的接口 (虽然这里是 Handler 结构体，非接口)
// 这是一个很好的实践，虽然对 Handler 本身不是必须的。

// TestValidateCreateRequest 单独测试校验逻辑
func TestValidateCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *cp_center.CreateCPMaterialRequest
		wantErr bool
		errText string
	}{
		{
			name: "Valid Request",
			req: &cp_center.CreateCPMaterialRequest{
				CPMaterial: &cp_center.CPMaterial{CpID: 1, CpName: "Valid Name"},
			},
			wantErr: false,
		},
		{
			name:    "Nil Request Material",
			req:     &cp_center.CreateCPMaterialRequest{CPMaterial: nil},
			wantErr: true,
			errText: "CPMaterial cannot be nil",
		},
		{
			name: "Missing CpID",
			req: &cp_center.CreateCPMaterialRequest{
				CPMaterial: &cp_center.CPMaterial{CpName: "Valid Name"},
			},
			wantErr: true,
			errText: "CpID is required",
		},
		{
			name: "Missing CpName",
			req: &cp_center.CreateCPMaterialRequest{
				CPMaterial: &cp_center.CPMaterial{CpID: 1, CpName: ""},
			},
			wantErr: true,
			errText: "CpName is required",
		},
		{
			name: "Whitespace CpName",
			req: &cp_center.CreateCPMaterialRequest{
				CPMaterial: &cp_center.CPMaterial{CpID: 1, CpName: "  \t "},
			},
			wantErr: true,
			errText: "CpName is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateRequest(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestNewMaterialFromRequest 单独测试模型转换
func TestNewMaterialFromRequest(t *testing.T) {
	setupIDGenerator() // 确保 idgen 工作

	reqBase := &cp_center.CreateCPMaterialRequest{
		CPMaterial: &cp_center.CPMaterial{
			CpID:               123,
			CpIcon:             "icon.png",
			CpName:             "Test CP",
			VerificationImages: []string{"img1.png", "img2.png"},
			BusinessLicenses:   "license.png",
			Website:            "test.com",
		},
	}

	t.Run("SubmitReview", func(t *testing.T) {
		req := *reqBase // 复制
		req.SubmitMode = cp_center.SubmitMode_SubmitReview

		mat, err := newMaterialFromRequest(&req)
		assert.NoError(t, err)
		assert.NotNil(t, mat)
		assert.Equal(t, uint64(123), mat.CpId)
		assert.Equal(t, "Test CP", mat.CpName)
		assert.Equal(t, `["img1.png","img2.png"]`, mat.VerificationImages)
		assert.Equal(t, int(cp_center.MaterialStatus_Reviewing), mat.Status)
		assert.True(t, mat.Id > 0)                                      // 确保 ID 已生成
		assert.WithinDuration(t, time.Now(), mat.CreateTs, time.Second) // 检查时间戳
	})

	t.Run("SubmitDraft", func(t *testing.T) {
		req := *reqBase // 复制
		req.SubmitMode = cp_center.SubmitMode_SubmitDraft

		mat, err := newMaterialFromRequest(&req)
		assert.NoError(t, err)
		assert.NotNil(t, mat)
		assert.Equal(t, int(cp_center.MaterialStatus_Draft), mat.Status)
	})

	t.Run("Default (SubmitDraft)", func(t *testing.T) {
		req := *reqBase     // 复制
		req.SubmitMode = 99 // 无效模式

		mat, err := newMaterialFromRequest(&req)
		assert.NoError(t, err)
		assert.NotNil(t, mat)
		assert.Equal(t, int(cp_center.MaterialStatus_Draft), mat.Status)
	})

	t.Run("Empty Images", func(t *testing.T) {
		req := *reqBase
		req.CPMaterial.VerificationImages = []string{} // 空数组
		req.SubmitMode = cp_center.SubmitMode_SubmitReview

		mat, err := newMaterialFromRequest(&req)
		assert.NoError(t, err)
		assert.NotNil(t, mat)
		assert.Equal(t, `[]`, mat.VerificationImages) // 应该为空 JSON 数组
	})

	t.Run("Nil Images", func(t *testing.T) {
		req := *reqBase
		req.CPMaterial.VerificationImages = nil // Nil slice
		req.SubmitMode = cp_center.SubmitMode_SubmitReview

		mat, err := newMaterialFromRequest(&req)
		assert.NoError(t, err)
		assert.NotNil(t, mat)
		assert.Equal(t, `[]`, mat.VerificationImages) // 也应该为空 JSON 数组
	})
}

// TestCheckCPExists 单独测试 CP 检查逻辑
func TestCheckCPExists(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCPRepo := mocks.NewMockICPRepo(mockCtrl)
	// 注意：checkCPExists 是 *CPMaterialHandler 的方法
	handler := &CPMaterialHandler{CPRepo: mockCPRepo}

	t.Run("Exists", func(t *testing.T) {
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(1)).
			Return(&ddl.GpCp{Id: 1}, nil)

		exists, err := handler.checkCPExists(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Not Exists", func(t *testing.T) {
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(2)).
			Return(nil, gorm.ErrRecordNotFound)

		exists, err := handler.checkCPExists(ctx, 2)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("DB Error", func(t *testing.T) {
		dbErr := errors.New("connection timeout")
		mockCPRepo.EXPECT().
			GetCPByID(ctx, int64(3)).
			Return(nil, dbErr)

		exists, err := handler.checkCPExists(ctx, 3)
		assert.Error(t, err)
		assert.False(t, exists)
		assert.True(t, errors.Is(err, dbErr))
	})
}
