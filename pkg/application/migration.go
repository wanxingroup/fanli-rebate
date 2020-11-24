package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/model/rebate"
)

func autoMigration() {
	db := database.GetDB(constant.DatabaseConfigKey)
	db.AutoMigrate(rebate.RebateItem{})
	db.AutoMigrate(rebate.RebateOrder{})
}
