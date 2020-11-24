package test

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/model/rebate"
)

var databaseName string

func Init() {

	rand.Seed(time.Now().UnixNano())

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("dir = %s\n", dir)
	for {
		dir = path.Dir(dir)
		if strings.LastIndex(dir, "/pkg") <= 0 {
			break
		}
	}

	fmt.Printf("processed dir = %s\n", dir)
	viper.SetConfigFile(fmt.Sprintf("%s/config/test.yaml", dir))
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	logrus.WithField("config", viper.AllSettings()).Info("viper loaded configuration")

	if err := viper.Unmarshal(config.Config); err != nil {
		panic(err)
	}

	logrus.WithField("config", config.Config).Info("unmarshalled configuration")

	databaseName = fmt.Sprintf("%s_%d", config.Config.MySQL[constant.DatabaseConfigKey].GetDatabase(), rand.Uint64())

	logLevel := config.Config.Log.GetLogLevel()
	logrus.SetLevel(logLevel)
	logrus.Debugf("log level set: %s", config.Config.Log.GetLogLevel())

	logrus.Infof("database name: %s", databaseName)

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&loc=Local&parseTime=true",
			config.Config.MySQL[constant.DatabaseConfigKey].GetUsername(),
			config.Config.MySQL[constant.DatabaseConfigKey].GetPassword(),
			config.Config.MySQL[constant.DatabaseConfigKey].GetHost(),
			config.Config.MySQL[constant.DatabaseConfigKey].GetPort(),
		))
	if err != nil {
		panic(err)
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` character set UTF8mb4 collate utf8mb4_general_ci", databaseName)).Error
	if err != nil {
		panic(err)
	}
	_ = db.Close()

	err = database.Connect(constant.DatabaseConfigKey, database.NewMySQLConfig(
		database.MySQLHost(config.Config.MySQL[constant.DatabaseConfigKey].GetHost()),
		database.MySQLPort(config.Config.MySQL[constant.DatabaseConfigKey].GetPort()),
		database.MySQLUsername(config.Config.MySQL[constant.DatabaseConfigKey].GetUsername()),
		database.MySQLPassword(config.Config.MySQL[constant.DatabaseConfigKey].GetPassword()),
		database.MySQLDatabase(databaseName),
	))

	if err != nil {
		panic(err)
	}

	objects := []interface{}{
		rebate.RebateItem{},
		rebate.RebateOrder{},
	}

	for _, object := range objects {
		err = database.GetDB(constant.DatabaseConfigKey).AutoMigrate(object).Error
		if err != nil {
			logrus.WithError(err).Error("migrate error")
		}
	}
}
