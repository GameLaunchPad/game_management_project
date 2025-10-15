package ddl

import (
	"time"
)

// 厂商信息
type GpCp struct {
	Id               uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;comment:cp_id" json:"id"`
	CpName           string    `gorm:"column:cp_name;type:varchar(512);comment:厂商名字;NOT NULL" json:"cp_name"`
	NewestMaterialId uint64    `gorm:"column:newest_material_id;type:bigint(20) unsigned;comment:最新资质id;NOT NULL" json:"newest_material_id"`
	OnlineMaterialId uint64    `gorm:"column:online_material_id;type:bigint(20) unsigned;comment:上线资质id;NOT NULL" json:"online_material_id"`
	VerifyStatus     uint      `gorm:"column:verify_status;type:int(11) unsigned;comment:0-未认证，1-已认证;NOT NULL" json:"verify_status"`
	CreateTs         time.Time `gorm:"column:create_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_ts"`
	ModifyTs         time.Time `gorm:"column:modify_ts;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"modify_ts"`
}

func (m *GpCp) TableName() string {
	return "gp_cp"
}
