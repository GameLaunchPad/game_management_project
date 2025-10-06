package handler

import (
	"encoding/json"
	"fmt"

	"github.com/GameLaunchPad/game_management_project/dao"
	"github.com/GameLaunchPad/game_management_project/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/kitex_gen/game"
)

func ConvertGameVersionToDdl(version *game.GameVersion) (*ddl.GpGameVersion, error) {
	if version == nil {
		return nil, fmt.Errorf("game version is nil")
	}

	platforms, err := json.Marshal(version.GamePlatforms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game platforms: %v", err)
	}
	images, err := json.Marshal(version.GameIntroductionImages)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game introduction images: %v", err)
	}

	return &ddl.GpGameVersion{
		GameName:               version.GameName,
		GameIcon:               version.GameIcon,
		HeaderImage:            version.HeaderImage,
		GameIntroduction:       version.GameIntroduction,
		GameIntroductionImages: string(images),
		Platform:               string(platforms),
		PackageName:            version.PackageName,
		DownloadUrl:            version.DownloadURL,
		Status:                 int(version.GameStatus),
	}, nil
}

func ConvertDdlToBriefGame(gameWithStatus *dao.GameWithVersionStatus) (*game.BriefGame, error) {
	if gameWithStatus == nil {
		return nil, fmt.Errorf("game is nil")
	}
	return &game.BriefGame{
		GameID:      int64(gameWithStatus.Id),
		CpID:        int64(gameWithStatus.CpId),
		GameName:    gameWithStatus.GameName,
		GameIcon:    gameWithStatus.GameIcon,
		HeaderImage: gameWithStatus.HeaderImage,
		CreateTime:  gameWithStatus.CreateTs.Unix(),
		UpdateTime:  gameWithStatus.ModifyTs.Unix(),
		GameStatus:  game.GameStatus(gameWithStatus.Status), // <--- 现在可以正确地从 JOIN 结果中获取 Status
	}, nil
}
