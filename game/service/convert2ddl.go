package service

import (
	"encoding/json"
	"fmt"

	"github.com/GameLaunchPad/game_management_project/game/dao"
	"github.com/GameLaunchPad/game_management_project/game/dao/ddl"
	"github.com/GameLaunchPad/game_management_project/game/kitex_gen/game"
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
		GameStatus:  game.GameStatus(gameWithStatus.Status),
	}, nil
}

// ConvertDdlToDetailGame converts GORM models to a GameDetail structure.
func ConvertDdlToDetailGame(gameDdl *ddl.GpGame, newestVersionDdl *ddl.GpGameVersion, onlineVersionDdl *ddl.GpGameVersion) (*game.GameDetail, error) {
	if gameDdl == nil {
		return nil, fmt.Errorf("input game ddl is nil")
	}

	newestVersion, err := ConvertDdlToGameVersion(newestVersionDdl)
	if err != nil {
		return nil, err
	}

	onlineVersion, err := ConvertDdlToGameVersion(onlineVersionDdl)
	if err != nil {
		return nil, err
	}

	return &game.GameDetail{
		GameID:             int64(gameDdl.Id),
		CpID:               int64(gameDdl.CpId),
		NewestGameVersion_: newestVersion,
		OnlineGameVersion:  onlineVersion,
		CreateTime:         gameDdl.CreateTs.Unix(),
		ModifyTime:         gameDdl.ModifyTs.Unix(),
	}, nil
}

// ConvertDdlToGameVersion converts a GORM model to a GameVersion structure.
func ConvertDdlToGameVersion(versionDdl *ddl.GpGameVersion) (*game.GameVersion, error) {
	if versionDdl == nil {
		return nil, nil
	}

	var platforms []game.GamePlatform
	if versionDdl.Platform != "" {
		if err := json.Unmarshal([]byte(versionDdl.Platform), &platforms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal platforms for version ID %d: %w", versionDdl.Id, err)
		}
	}

	var images []string
	if versionDdl.GameIntroductionImages != "" {
		if err := json.Unmarshal([]byte(versionDdl.GameIntroductionImages), &images); err != nil {
			return nil, fmt.Errorf("failed to unmarshal images for version ID %d: %w", versionDdl.Id, err)
		}
	}

	return &game.GameVersion{
		GameID:                 int64(versionDdl.GameId),
		GamVersionID:           int64(versionDdl.Id),
		GameName:               versionDdl.GameName,
		GameIcon:               versionDdl.GameIcon,
		HeaderImage:            versionDdl.HeaderImage,
		GameIntroduction:       versionDdl.GameIntroduction,
		GameIntroductionImages: images,
		GamePlatforms:          platforms,
		PackageName:            versionDdl.PackageName,
		DownloadURL:            versionDdl.DownloadUrl,
		GameStatus:             game.GameStatus(versionDdl.Status),
		ReviewComment:          versionDdl.ReviewComment,
		ReviewTime:             versionDdl.ReviewTime,
		CreateTime:             versionDdl.CreateTs.Unix(),
		UpdateTime:             versionDdl.ModifyTs.Unix(),
	}, nil
}
