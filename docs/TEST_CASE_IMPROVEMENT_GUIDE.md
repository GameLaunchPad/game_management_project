# 单元测试用例改进指南

本文档列出了需要添加的测试用例，以将覆盖率从 85% 提升到 95%。

## CreateGameDetail 函数

### 当前已有测试用例：
- ✅ `TestCreateGameDetail_Success` - 成功创建游戏
- ✅ `TestCreateGameDetail_FailWithNonZeroGameID` - GameID 不为 0 的情况
- ✅ `TestCreateGameDetail_DaoError` - DAO 返回错误

### 缺失的测试用例：

#### 1. TestCreateGameDetail_NilGameDetail
**代码位置**: `create_game_detail.go:17-23`
**测试场景**: GameDetail 为 nil
```go
func TestCreateGameDetail_NilGameDetail(t *testing.T) {
    req := &game.CreateGameDetailRequest{
        GameDetail: nil,
    }
    
    resp, err := CreateGameDetail(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "GameDetail or GameVersion is missing")
}
```

#### 2. TestCreateGameDetail_NilGameVersion
**代码位置**: `create_game_detail.go:17-23`
**测试场景**: GameVersion 为 nil
```go
func TestCreateGameDetail_NilGameVersion(t *testing.T) {
    req := &game.CreateGameDetailRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: 0,
            CpID:   1001,
            GameVersion: nil,
        },
    }
    
    resp, err := CreateGameDetail(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "GameDetail or GameVersion is missing")
}
```

#### 3. TestCreateGameDetail_ConvertError
**代码位置**: `create_game_detail.go:40-45`
**测试场景**: ConvertGameVersionToDdl 返回错误
```go
func TestCreateGameDetail_ConvertError(t *testing.T) {
    setupIDGenerator()
    
    // 注意：这个测试需要 mock service.ConvertGameVersionToDdl
    // 如果无法直接 mock，可能需要创建一个无效的 GameVersion 数据
    // 具体取决于 ConvertGameVersionToDdl 的实现
    
    req := &game.CreateGameDetailRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: 0,
            CpID:   1001,
            GameVersion: &game.GameVersion{
                // 创建一个会导致 ConvertGameVersionToDdl 失败的数据
                // 需要根据实际实现调整
            },
        },
    }
    
    resp, err := CreateGameDetail(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "Invalid game version data")
}
```

## UpdateGameDraft 函数

### 当前已有测试用例：
- ✅ `TestUpdateGameDraft_Success` - 成功更新草稿
- ✅ `TestUpdateGameDraft_GameNotFound` - 游戏不存在
- ✅ `TestUpdateGameDraft_FailWithZeroGameID` - GameID 为 0

### 缺失的测试用例：

#### 1. TestUpdateGameDraft_NilGameDetail
**代码位置**: `update_game_draft.go:16-20`
**测试场景**: GameDetail 为 nil
```go
func TestUpdateGameDraft_NilGameDetail(t *testing.T) {
    req := &game.UpdateGameDraftRequest{
        GameDetail: nil,
    }
    
    resp, err := UpdateGameDraft(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "GameDetail or GameVersion is missing")
}
```

#### 2. TestUpdateGameDraft_NilGameVersion
**代码位置**: `update_game_draft.go:16-20`
**测试场景**: GameVersion 为 nil
```go
func TestUpdateGameDraft_NilGameVersion(t *testing.T) {
    req := &game.UpdateGameDraftRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: 12345,
            CpID:   1001,
            GameVersion: nil,
        },
    }
    
    resp, err := UpdateGameDraft(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "GameDetail or GameVersion is missing")
}
```

#### 3. TestUpdateGameDraft_NegativeGameID
**代码位置**: `update_game_draft.go:22-26`
**测试场景**: GameID 为负数
```go
func TestUpdateGameDraft_NegativeGameID(t *testing.T) {
    req := &game.UpdateGameDraftRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: -1,
            CpID:   1001,
            GameVersion: &game.GameVersion{
                GameName: "Test Game",
            },
        },
    }
    
    resp, err := UpdateGameDraft(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "GameID is required")
}
```

#### 4. TestUpdateGameDraft_ConvertError
**代码位置**: `update_game_draft.go:33-38`
**测试场景**: ConvertGameVersionToDdl 返回错误
```go
func TestUpdateGameDraft_ConvertError(t *testing.T) {
    setupIDGenerator()
    
    // 创建一个会导致 ConvertGameVersionToDdl 失败的数据
    req := &game.UpdateGameDraftRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: 12345,
            CpID:   1001,
            GameVersion: &game.GameVersion{
                // 无效数据，需要根据实际实现调整
            },
        },
    }
    
    resp, err := UpdateGameDraft(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "Invalid game version data")
}
```

#### 5. TestUpdateGameDraft_OtherError
**代码位置**: `update_game_draft.go:46-55`
**测试场景**: DAO 返回其他错误（非 RecordNotFound）
```go
func TestUpdateGameDraft_OtherError(t *testing.T) {
    setupIDGenerator()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    otherError := errors.New("database connection error")
    mockGameDAO.EXPECT().
        UpdateGameDraft(gomock.Any(), gomock.Any(), gomock.Any()).
        Return(otherError).
        Times(1)
    
    req := &game.UpdateGameDraftRequest{
        GameDetail: &game.GameDetailWrite{
            GameID: 12345,
            CpID:   1001,
            GameVersion: &game.GameVersion{
                GameName: "Test Game",
            },
        },
    }
    
    resp, err := UpdateGameDraft(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "500", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "Internal Server Error")
}
```

## GetGameList 函数

### 当前已有测试用例：
- ✅ `TestGetGameList_Success` - 成功获取列表
- ✅ `TestGetGameList_WithFilter` - 带过滤条件
- ✅ `TestGetGameList_DaoError` - DAO 返回错误

### 缺失的测试用例：

#### 1. TestGetGameList_InvalidPageNum
**代码位置**: `get_game_list.go:19-22`
**测试场景**: PageNum 为 0 或负数（应该默认设置为 1）
```go
func TestGetGameList_InvalidPageNum(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    // 期望使用默认值 1
    mockGameDAO.EXPECT().
        GetGameList(gomock.Any(), gomock.Any(), 1, 10).
        Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
        Times(1)
    
    req := &game.GetGameListRequest{
        PageNum:  0,  // 无效值，应该使用默认值 1
        PageSize: 10,
    }
    
    resp, err := GetGameList(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "200", resp.BaseResp.Code)
}
```

#### 2. TestGetGameList_InvalidPageSize
**代码位置**: `get_game_list.go:24-27`
**测试场景**: PageSize 为 0 或负数（应该默认设置为 10）
```go
func TestGetGameList_InvalidPageSize(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    // 期望使用默认值 10
    mockGameDAO.EXPECT().
        GetGameList(gomock.Any(), gomock.Any(), 1, 10).
        Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
        Times(1)
    
    req := &game.GetGameListRequest{
        PageNum:  1,
        PageSize: 0,  // 无效值，应该使用默认值 10
    }
    
    resp, err := GetGameList(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "200", resp.BaseResp.Code)
}
```

#### 3. TestGetGameList_NoFilter
**代码位置**: `get_game_list.go:14-17`
**测试场景**: Filter 未设置或为空（应该使用 nil）
```go
func TestGetGameList_NoFilter(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    // 期望 filterText 为 nil
    mockGameDAO.EXPECT().
        GetGameList(gomock.Any(), nil, 1, 10).
        Return([]*dao.GameWithVersionStatus{}, int64(0), nil).
        Times(1)
    
    req := &game.GetGameListRequest{
        PageNum:  1,
        PageSize: 10,
        // Filter 未设置
    }
    
    resp, err := GetGameList(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "200", resp.BaseResp.Code)
}
```

#### 4. TestGetGameList_ConvertError
**代码位置**: `get_game_list.go:40-45`
**测试场景**: ConvertDdlToBriefGame 返回错误
```go
func TestGetGameList_ConvertError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    // 返回一个会导致转换失败的数据
    // 需要根据 ConvertDdlToBriefGame 的实际实现调整
    invalidGame := &dao.GameWithVersionStatus{
        GpGame: ddl.GpGame{
            // 创建会导致转换失败的数据
        },
    }
    
    mockGameDAO.EXPECT().
        GetGameList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
        Return([]*dao.GameWithVersionStatus{invalidGame}, int64(1), nil).
        Times(1)
    
    req := &game.GetGameListRequest{
        PageNum:  1,
        PageSize: 10,
    }
    
    resp, err := GetGameList(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "500", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "Failed to convert game data")
}
```

## ReviewGameVersion 函数

### 当前已有测试用例：
- ✅ `TestReviewGameVersion_PassSuccess` - 审核通过
- ✅ `TestReviewGameVersion_RejectSuccess` - 审核拒绝
- ✅ `TestReviewGameVersion_NotFound` - 游戏不存在

### 缺失的测试用例：

#### 1. TestReviewGameVersion_InvalidGameID
**代码位置**: `review_game_version.go:14-18`
**测试场景**: GameID 为 0 或负数
```go
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
```

#### 2. TestReviewGameVersion_InvalidVersionID
**代码位置**: `review_game_version.go:14-18`
**测试场景**: GameVersionID 为 0 或负数
```go
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
```

#### 3. TestReviewGameVersion_InvalidReviewResult
**代码位置**: `review_game_version.go:27-31`
**测试场景**: ReviewResult 为无效值
```go
func TestReviewGameVersion_InvalidReviewResult(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockGameDAO := mock.NewMockIGameDAO(ctrl)
    GameDao = mockGameDAO
    
    // 使用一个无效的 ReviewResult 值
    // 注意：需要根据实际的枚举值调整
    req := &game.ReviewGameVersionRequest{
        GameID:        101,
        GameVersionID: 201,
        ReviewResult_: game.ReviewResult__(999), // 无效值
    }
    
    resp, err := ReviewGameVersion(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "400", resp.BaseResp.Code)
    assert.Contains(t, resp.BaseResp.Msg, "Invalid review result")
}
```

#### 4. TestReviewGameVersion_OtherError
**代码位置**: `review_game_version.go:46-50`
**测试场景**: DAO 返回其他错误（非 RecordNotFound）
```go
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
```

## GetGameDetail 函数

### 当前已有测试用例：
- ✅ `TestGetGameDetail_Success` - 成功获取详情
- ✅ `TestGetGameDetail_GameNotFound` - 游戏不存在
- ✅ `TestGetGameDetail_DaoError` - DAO 返回错误

### 可能需要添加的测试用例：

#### 1. TestGetGameDetail_InvalidGameID
**测试场景**: GameID 为 0 或负数（如果代码中有验证）

## 测试用例添加优先级

### 高优先级（直接影响覆盖率）：
1. ✅ CreateGameDetail - NilGameDetail
2. ✅ CreateGameDetail - NilGameVersion
3. ✅ UpdateGameDraft - NilGameDetail
4. ✅ UpdateGameDraft - NilGameVersion
5. ✅ UpdateGameDraft - NegativeGameID
6. ✅ ReviewGameVersion - InvalidGameID
7. ✅ ReviewGameVersion - InvalidVersionID
8. ✅ GetGameList - InvalidPageNum
9. ✅ GetGameList - InvalidPageSize

### 中优先级（需要 mock 服务函数）：
1. CreateGameDetail - ConvertError
2. UpdateGameDraft - ConvertError
3. GetGameList - ConvertError

### 低优先级（边界情况）：
1. UpdateGameDraft - OtherError
2. ReviewGameVersion - InvalidReviewResult
3. ReviewGameVersion - OtherError

## 添加测试用例的步骤

1. **复制测试用例模板**：从上面的代码示例中复制对应的测试用例
2. **调整导入**：确保导入必要的包（gomock, assert 等）
3. **运行测试**：确保测试能够编译和运行
4. **验证覆盖率**：查看HTML报告 `game/handler_coverage.html`（已包含总覆盖率信息）
5. **检查覆盖率报告**：打开 `game/handler_coverage.html` 确认未覆盖的代码已被覆盖

## 运行测试并查看覆盖率

```bash
# 运行测试
bash scripts/run_tests_simple.sh

# 查看覆盖率分析
# 查看HTML覆盖率报告（已包含总覆盖率信息）
open game/handler_coverage.html

# 查看HTML覆盖率报告
open game/handler_coverage.html
```

## 注意事项

1. **Mock 服务函数**：如果测试需要 mock `service.ConvertGameVersionToDdl` 或 `service.ConvertDdlToBriefGame`，可能需要：
   - 创建 service 包的 mock
   - 或者创建会导致转换失败的实际数据

2. **测试数据**：确保测试数据符合业务逻辑要求

3. **错误消息**：确保断言中的错误消息与实际代码中的消息匹配

4. **代码覆盖率**：添加测试用例后，运行覆盖率分析，确认覆盖率提升

