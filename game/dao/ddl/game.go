package ddl

import "time"

// 游戏信息
type GpGame struct {
	Id                     uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;comment:游戏ID" json:"id"`
	CpId                   uint64    `gorm:"column:cp_id;type:bigint(20) unsigned;comment:厂商ID;NOT NULL" json:"cp_id"`
	GameName               string    `gorm:"column:game_name;type:varchar(1024);comment:游戏名;NOT NULL" json:"game_name"`
	GameIcon               string    `gorm:"column:game_icon;type:varchar(512);comment:游戏图片URI;NOT NULL" json:"game_icon"`
	HeaderImage            string    `gorm:"column:header_image;type:varchar(512);comment:游戏头图URI;NOT NULL" json:"header_image"`
	GameIntroduction       string    `gorm:"column:game_introduction;type:text;comment:游戏简介" json:"game_introduction"`
	GameIntroductionImages string    `gorm:"column:game_introduction_images;type:text;comment:游戏介绍图" json:"game_introduction_images"`
	Platform               string    `gorm:"column:platform;type:varchar(256);comment:游戏推广平台 0-unset, 1-android, 2-ios, 3-web,可以支持多端配置，为Json数组;NOT NULL" json:"platform"`
	PackageName            string    `gorm:"column:package_name;type:varchar(256);comment:游戏包名（APP端使用）;NOT NULL" json:"package_name"`
	DownloadUrl            string    `gorm:"column:download_url;type:text;comment:游戏下载链接" json:"download_url"`
	NewestGameVersionId    uint64    `gorm:"column:newest_game_version_id;type:bigint(20) unsigned;comment:最新游戏版本id" json:"newest_game_version_id"`
	OnlineGameVersionId    uint64    `gorm:"column:online_game_version_id;type:bigint(20) unsigned;comment:上线游戏版本id" json:"online_game_version_id"`
	CreateTs               time.Time `gorm:"column:create_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_ts"`
	ModifyTs               time.Time `gorm:"column:modify_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"modify_ts"`
}

func (m *GpGame) TableName() string {
	return "gp_game"
}
