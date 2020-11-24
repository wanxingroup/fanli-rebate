package test

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
)

func Release() {

	var err error
	err = database.GetDB(constant.DatabaseConfigKey).Exec(fmt.Sprintf("DROP DATABASE %s", databaseName)).Error
	if err != nil {
		logrus.WithField("error", err).Error("drop database error")
	}

	err = database.Disconnect(constant.DatabaseConfigKey)
	if err != nil {
		logrus.WithField("error", err).Error("close database error")
	}
}
