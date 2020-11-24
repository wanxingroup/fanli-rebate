package application

import (
	"github.com/spf13/viper"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/rebate"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func Start() {

	app := launcher.NewApplication(
		launcher.SetApplicationDescription(
			&launcher.ApplicationDescription{
				ShortDescription: "rebate service",
				LongDescription:  "rebate",
			},
		),
		launcher.SetApplicationLogger(log.GetLogger()),
		launcher.SetApplicationEvents(
			launcher.NewApplicationEvents(
				launcher.SetOnInitEvent(func(app *launcher.Application) {

					unmarshalConfiguration()
					registerRebateRPCRouter(app)

					client.InitUserService()
					client.InitMerchantService()
					client.InitCardService()

				}),

				launcher.SetOnStartEvent(func(app *launcher.Application) {

					autoMigration()
				}),
			),
		),
	)

	app.Launch()
}
func registerRebateRPCRouter(app *launcher.Application) {

	rpcService := app.GetRPCService()
	if rpcService == nil {

		log.GetLogger().WithField("stage", "onInit").Error("get rpc service is nil")
		return
	}

	protos.RegisterRebateControllerServer(rpcService.GetRPCConnection(), &rebate.Controller{})
}
func unmarshalConfiguration() {
	err := viper.Unmarshal(config.Config)
	if err != nil {

		log.GetLogger().WithError(err).Error("unmarshal config error")
	}

}
