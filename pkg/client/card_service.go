package client

import (
	"context"

	cardProtos "dev-gitlab.wanxingrowth.com/fanli/card/pkg/rpc/protos"
	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

var cardRPCService cardProtos.CardControllerClient

func InitCardService() {

	log.GetLogger().Info("starting init card rpc service")

	var ctx = context.Background()

	var rpcConfig, exist = config.Config.RPCServices[constant.RPCCardServiceConfigKey]
	if !exist {
		log.GetLogger().Error("card rpc service configuration not exist")
		return
	}

	if rpcConfig.GetConnectTimeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.TODO(), rpcConfig.GetConnectTimeout())
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, rpcConfig.GetAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.GetLogger().WithError(err).Error("card rpc service connect failed")
		return
	}

	cardRPCService = cardProtos.NewCardControllerClient(conn)

	log.GetLogger().Info("card rpc service init succeed")
}

func GetCardService() cardProtos.CardControllerClient {
	return cardRPCService
}
