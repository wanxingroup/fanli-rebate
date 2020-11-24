package rebate

import (
	database "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameRebateItem = "rebate_item"

type RebateItem struct {
	ItemId            uint64 `gorm:"column:itemId;type:bigint unsigned;primary_key;comment:'商品 Id'"`
	ItemType          uint8  `gorm:"column:itemType;type:tinyint unsigned;primary_key;comment:'类型 1 商品 2 权益卡'"`
	FirstRebateRatio  uint32 `gorm:"column:firstRebateRatio;type:tinyint unsigned;comment:'一级拨比'"`
	SecondRebateRatio uint32 `gorm:"column:secondRebateRatio;type:tinyint unsigned;comment:'二级拨比'"`
	CreatedBy         uint64 `gorm:"column:createdBy;type:bigint unsigned;not null;default: '0';comment:'创建人'"`
	UpdatedBy         uint64 `gorm:"column:updatedBy;type:bigint unsigned;not null;default: '0';comment:'更新人'"`
	database.Time
}

func (rebate *RebateItem) TableName() string {
	return TableNameRebateItem
}
