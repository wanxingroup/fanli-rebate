package rebate

import (
	"time"

	database "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameRebateOrder = "rebate_order"

type RebateOrder struct {
	OrderId          uint64    `gorm:"column:orderId;type:bigint unsigned;primary_key;comment:'order Id'"`
	ShopId           uint64    `gorm:"column:shopId;type:bigint unsigned;primary_key;comment:'shop Id'"`
	UserId           uint64    `gorm:"column:userId;type:bigint unsigned;not null;index:userId;comment:'用户 Id'"`
	Mobile           string    `gorm:"column:mobile;type:char(11);not null;index:mobile;comment:'用户 手机号'"`
	PaidPrice        uint64    `gorm:"column:paidPrice;type:bigint unsigned;not null;comment:'实际支付金额'"`
	FirstRebate      uint64    `gorm:"column:firstRebate;type:bigint unsigned;comment:'一级收益'"`
	FirstRebateRate  uint64    `gorm:"column:firstRebateRate;type:bigint unsigned;comment:'一级收益的百分比'"`
	FirstUserId      uint64    `gorm:"column:firstUserId;type:bigint unsigned;comment:'直接上级的 userId'"`
	FirstType        uint8     `gorm:"column:firstType;type:tinyint unsigned;comment:'一级分类的区别 1: 平台 2：加盟商 3：门店 4：店员 5：用户'"`
	SecondRebate     uint64    `gorm:"column:secondRebate;type:bigint unsigned;comment:'二级受益'"`
	SecondRebateRate uint64    `gorm:"column:secondRebateRate;type:bigint unsigned;comment:'二级受益百分比'"`
	SecondUserId     uint64    `gorm:"column:secondUserId;type:bigint unsigned;comment:'间接上级的 userId'"`
	SecondType       uint8     `gorm:"column:secondType;type:tinyint unsigned;comment:'二级分类的区别 1: 平台 2：加盟商 3：门店 4：店员 5：用户'"`
	PaidTime         time.Time `gorm:"column:paidTime;null;comment:'支付时间'"`
	ItemType         uint8     `gorm:"column:itemType;type:tinyint unsigned;comment:'类型 1 商品 2 权益卡'"`
	ItemId           uint64    `gorm:"column:itemId;type:bigint unsigned;comment:'商品 Id'"`
	database.Time
}

func (rebate *RebateOrder) TableName() string {
	return TableNameRebateOrder
}

const (
	InviterTypeDefault    = 1 //默认平台
	InviterTypeFranchisee = 2 //加盟商
	InviterTypeShop       = 3 //店铺
	InviterTypeStaff      = 4 //员工
	InviterTypeUser       = 5 //用户

	ItemTypeGoods = 1 // 订单类型为 商品
	ItemTypeCard  = 2 // 订单类型为 权益卡
)

type UserInvited struct {
	UserId     uint64 // 用户自己的Id
	UserMobile string // 用户手机号
	UserType   uint32 // 自己类型
	UserName   string // 用户名字
	IsHaveCard bool   // 是否拥有权益卡
	IsRebate   bool   // 是否返利

	InviterId       uint64 // 邀请者 Id
	InviterType     uint32 // 邀请者类型
	InviterIsRebate bool   // 邀请者是否开启返利
}

// 统计需要的返回 数据结构
type StatisticsOrderResult struct {
	PayMoneyTotal     uint64 `gorm:"column:payMoneyTotal"`     // 总支付金额
	FirstRebateTotal  uint64 `gorm:"column:firstRebateTotal"`  // 一级返利
	SecondRebateTotal uint64 `gorm:"column:secondRebateTotal"` // 二级返利
}
