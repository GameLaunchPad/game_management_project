package ddl

import (
	"time"
)

// 资质材料信息
type GpCpMaterial struct {
	Id                 uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;comment:材料ID" json:"id"`
	CpId               uint64    `gorm:"column:cp_id;type:bigint(20) unsigned;comment:厂商id;NOT NULL" json:"cp_id"`
	CpIcon             string    `gorm:"column:cp_icon;type:varchar(256);comment:厂商ICON;NOT NULL" json:"cp_icon"`
	CpName             string    `gorm:"column:cp_name;type:varchar(512);comment:厂商名字;NOT NULL" json:"cp_name"`
	VerificationImages string    `gorm:"column:verification_images;type:text;comment:资质图片URI（Json）" json:"verification_images"`
	BusinessLicense    string    `gorm:"column:business_license;type:varchar(2048);comment:营业执照;NOT NULL" json:"business_license"`
	Website            string    `gorm:"column:website;type:text;comment:官网地址" json:"website"`
	Status             int       `gorm:"column:status;type:int(11);comment:0-Unset, 1-草稿, 2-审核中, 3-已发布，4-已拒绝;NOT NULL" json:"status"`
	Operator           string    `gorm:"column:operator;type:varchar(128);comment:审核人;NOT NULL" json:"operator"`
	ReviewComment      string    `gorm:"column:review_comment;type:text;comment:审核意见" json:"review_comment"`
	CreateTs           time.Time `gorm:"column:create_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_ts"`
	ModifyTs           time.Time `gorm:"column:modify_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"modify_ts"`
}

func (m *GpCpMaterial) TableName() string {
	return "gp_cp_material"
}
