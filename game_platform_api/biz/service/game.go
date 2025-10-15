package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/biz/model/game_platform_api"
	"github.com/GameLaunchPad/game_management_project/game_platform_api/rpc"
)

// GameService 封装了所有与下游 game 服务相关的业务逻辑
type GameService struct{}

// NewGameService 创建一个新的 GameService 实例
func NewGameService() *GameService {
	return &GameService{}
}

// GetGameList 调用 game 服务获取游戏列表
func (s *GameService) GetGameList(ctx context.Context, req *game_platform_api.GetGameListRequest) (*game.GetGameListResponse, error) {
	rpcReq := &game.GetGameListRequest{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	if req.Filter != nil {
		rpcReq.Filter = &game.GameListFilter{
			FilterText: req.Filter.FilterText,
		}
	}
	if req.Sorter != nil {
		rpcReq.Sorter = &game.GameListSorter{
			UpdateTime: req.Sorter.UpdateTime,
		}
	}

	resp, err := rpc.GameClient.GetGameList(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetGameDetail 调用 game 服务获取游戏详情
func (s *GameService) GetGameDetail(ctx context.Context, req *game_platform_api.GetGameDetailRequest) (*game.GetGameDetailResponse, error) {
	rpcReq := &game.GetGameDetailRequest{
		GameID: req.GameID,
	}
	resp, err := rpc.GameClient.GetGameDetail(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateGameDetail 调用 game 服务创建游戏详情
func (s *GameService) CreateGameDetail(ctx context.Context, req *game_platform_api.CreateGameDetailRequest) (*game.CreateGameDetailResponse, error) {
	rpcReq := &game.CreateGameDetailRequest{
		GameDetail: &game.GameDetailWrite{
			CpID: req.GameDetail.CpID,
			GameVersion: &game.GameVersion{
				GameName:               req.GameDetail.GameVersion.GameName,
				GameIcon:               req.GameDetail.GameVersion.GameIcon,
				HeaderImage:            req.GameDetail.GameVersion.HeaderImage,
				GameIntroduction:       req.GameDetail.GameVersion.GameIntroduction,
				GameIntroductionImages: req.GameDetail.GameVersion.GameIntroductionImages,
				GamePlatforms:          convertPlatformToRPC(req.GameDetail.GameVersion.GamePlatforms),
				PackageName:            req.GameDetail.GameVersion.PackageName,
				DownloadURL:            req.GameDetail.GameVersion.DownloadURL,
			},
		},
		SubmitMode: convertSubmitModeToRPC(req.SubmitMode),
	}

	resp, err := rpc.GameClient.CreateGameDetail(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateGameDetail 调用 game 服务更新游戏详情
func (s *GameService) UpdateGameDetail(ctx context.Context, req *game_platform_api.UpdateGameDetailRequest) (*game.UpdateGameDraftResponse, error) {
	gameID, err := strconv.ParseInt(req.GameID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid game_id format: %w", err)
	}

	gameIDInDetail, err := strconv.ParseInt(req.GameDetail.GameID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid game_detail.game_id format: %w", err)
	}

	rpcReq := &game.UpdateGameDraftRequest{
		GameDetail: &game.GameDetailWrite{
			GameID: gameIDInDetail,
			CpID:   req.GameDetail.CpID,
			GameVersion: &game.GameVersion{
				GameID:                 gameID, // 确保版本信息中的 GameID 也被设置
				GameName:               req.GameDetail.GameVersion.GameName,
				GameIcon:               req.GameDetail.GameVersion.GameIcon,
				HeaderImage:            req.GameDetail.GameVersion.HeaderImage,
				GameIntroduction:       req.GameDetail.GameVersion.GameIntroduction,
				GameIntroductionImages: req.GameDetail.GameVersion.GameIntroductionImages,
				GamePlatforms:          convertPlatformToRPC(req.GameDetail.GameVersion.GamePlatforms),
				PackageName:            req.GameDetail.GameVersion.PackageName,
				DownloadURL:            req.GameDetail.GameVersion.DownloadURL,
			},
		},
		SubmitMode: convertSubmitModeToRPC(req.SubmitMode),
	}

	resp, err := rpc.GameClient.UpdateGameDraft(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ReviewGameVersion 调用 game 服务审核游戏版本
func (s *GameService) ReviewGameVersion(ctx context.Context, req *game_platform_api.ReviewGameVersionRequest) (*game.ReviewGameVersionResponse, error) {
	gameID, err := strconv.ParseInt(req.GameID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid game_id format: %w", err)
	}

	gameVersionID, err := strconv.ParseInt(req.GameVersionID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid game_version_id format: %w", err)
	}

	rpcReq := &game.ReviewGameVersionRequest{
		GameID:        gameID,
		GameVersionID: gameVersionID,
		ReviewResult_: convertReviewResultToRPC(req.ReviewResult),
	}

	resp, err := rpc.GameClient.ReviewGameVersion(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteGameDraft calls the game service to delete a game draft.
func (s *GameService) DeleteGameDraft(ctx context.Context, req *game_platform_api.DeleteGameDraftRequest) (*game.DeleteGameDraftResponse, error) {
	rpcReq := &game.DeleteGameDraftRequest{
		GameID: req.GameID, // 修正: GameID
	}
	resp, err := rpc.GameClient.DeleteGameDraft(ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- 类型转换辅助函数 ---

func convertSubmitModeToRPC(mode game_platform_api.SubmitMode) game.SubmitMode {
	switch mode {
	case game_platform_api.SubmitMode_SubmitDraft:
		return game.SubmitMode_SubmitDraft
	case game_platform_api.SubmitMode_SubmitReview:
		return game.SubmitMode_SubmitReview
	default:
		return game.SubmitMode_Unset
	}
}

// 修正：确保函数名和返回值与 game.thrift 定义一致
func convertReviewResultToRPC(result game_platform_api.ReviewResult) game.ReviewResult_ {
	switch result {
	case game_platform_api.ReviewResult_Pass:
		return game.ReviewResult__Pass
	case game_platform_api.ReviewResult_Reject:
		return game.ReviewResult__Reject
	default:
		return game.ReviewResult__Unset
	}
}

func convertPlatformToRPC(platforms []game_platform_api.GamePlatform) []game.GamePlatform {
	rpcPlatforms := make([]game.GamePlatform, 0, len(platforms))
	for _, p := range platforms {
		switch p {
		case game_platform_api.GamePlatform_Android:
			rpcPlatforms = append(rpcPlatforms, game.GamePlatform_Android)
		case game_platform_api.GamePlatform_IOS:
			rpcPlatforms = append(rpcPlatforms, game.GamePlatform_IOS)
		case game_platform_api.GamePlatform_Web:
			rpcPlatforms = append(rpcPlatforms, game.GamePlatform_Web)
		}
	}
	return rpcPlatforms
}
