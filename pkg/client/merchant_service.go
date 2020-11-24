package client

import (
	"context"

	merchantProtos "dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

var merchantRPCService merchantProtos.MerchantControllerClient
var shopRPCService merchantProtos.ShopControllerClient
var staffRPCService merchantProtos.StaffControllerClient

func InitMerchantService() {

	log.GetLogger().Info("starting init merchant rpc service")

	var ctx = context.Background()

	var rpcConfig, exist = config.Config.RPCServices[constant.RPCMerchantServiceConfigKey]
	if !exist {
		log.GetLogger().Error("merchant rpc service configuration not exist")
		return
	}

	if rpcConfig.GetConnectTimeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.TODO(), rpcConfig.GetConnectTimeout())
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, rpcConfig.GetAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.GetLogger().WithError(err).Error("merchant rpc service connect failed")
		return
	}

	merchantRPCService = merchantProtos.NewMerchantControllerClient(conn)
	shopRPCService = merchantProtos.NewShopControllerClient(conn)
	staffRPCService = merchantProtos.NewStaffControllerClient(conn)

	log.GetLogger().Info("merchant rpc service init succeed")
}

func GetMerchantService() merchantProtos.MerchantControllerClient {
	return merchantRPCService
}

func GetShopService() merchantProtos.ShopControllerClient {
	return shopRPCService
}

func GetStaffService() merchantProtos.StaffControllerClient {
	return staffRPCService
}
