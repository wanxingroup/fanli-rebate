package rebate

import (
	"fmt"

	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/model/rebate"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func (_ Controller) CreateRebateOrder(ctx context.Context, req *protos.CreateRebateOrderRequest) (*protos.CreateRebateOrderReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}
	var recodeData = &protos.RebateOrder{
		OrderId:   req.RebateOrder.OrderId,
		PaidPrice: req.RebateOrder.PaidPrice,
		PaidTime:  req.RebateOrder.PaidTime,
		ItemType:  req.RebateOrder.ItemType,
		ShopId:    req.RebateOrder.ShopId,
		ItemId:    req.RebateOrder.ItemId,
	}

	logger.WithField("recoodeData", recodeData).Info("recoodeData")
	// 获取该商品的一二级拨比比例

	// 获取原来的数据
	rebateItem, rebateItemErr := GetRebateByIteam(req.RebateOrder.ItemId, uint8(req.RebateOrder.ItemType))
	if rebateItemErr != nil {
		logger.WithError(rebateItemErr).Error("get rebate error")
		// 默认拨比给 0 0
		recodeData.FirstRebateRate = 0
		recodeData.SecondRebateRate = 0
	} else {
		recodeData.FirstRebateRate = uint64(rebateItem.FirstRebateRatio)
		recodeData.SecondRebateRate = uint64(rebateItem.SecondRebateRatio)
	}

	// 获取当前用户的信息
	currentUserInfo, err := GetInviterInfo(req.RebateOrder.UserId, req.RebateOrder.ShopId, rebate.InviterTypeUser, logger)

	if err != nil {
		logger.WithField("err", err).Error("get currentUserInfo error")
		return &protos.CreateRebateOrderReply{
			Err: &protos.Error{
				Code:    err.Code,
				Message: err.Message,
				Stack:   nil,
			},
		}, nil
	}

	// 获取当前用户是否有权益卡, 如果没有直接返回就可以了
	if !GetUserIsHaveCard(logger, req.RebateOrder.UserId) {
		return &protos.CreateRebateOrderReply{
			RebateOrder: nil,
		}, nil
	}

	logger.WithField("currentUserInfo", currentUserInfo).Info("currentUserInfo")

	recodeData.UserId = currentUserInfo.UserId
	recodeData.Mobile = currentUserInfo.UserMobile

	/*** 一级返利 ***/
	recodeData.FirstType = currentUserInfo.InviterType
	recodeData.FirstUserId = currentUserInfo.InviterId
	recodeData.FirstRebate = recodeData.FirstRebateRate * recodeData.PaidPrice / 100
	// 获取一级返利者信息
	var firstUserInfo *rebate.UserInvited
	firstUserInfo, err = GetInviterInfo(currentUserInfo.InviterId, req.RebateOrder.ShopId, currentUserInfo.InviterType, logger)

	if err != nil {
		logger.WithField("err", err).Error("get firstUserInfo error")
		return &protos.CreateRebateOrderReply{
			Err: &protos.Error{
				Code:    err.Code,
				Message: err.Message,
				Stack:   nil,
			},
		}, nil
	}

	/**
	 * 1.1 判断一级是否开启返利, 如果关闭返利以及是员工/门店, 收益归加盟商
	 */
	if !firstUserInfo.IsRebate {
		// 如果是门店, 收益归加盟商
		if firstUserInfo.UserType == rebate.InviterTypeShop {
			recodeData.FirstType = firstUserInfo.InviterType
			recodeData.FirstUserId = firstUserInfo.InviterId
		}
		// 如果是员工, 收益也归加盟商
		if firstUserInfo.UserType == rebate.InviterTypeStaff {
			// 通过员工找上一级的Id(门店)找加盟商
			shopInfo, shopInfoErr := GetInviterInfo(firstUserInfo.InviterId, req.RebateOrder.ShopId, firstUserInfo.InviterType, logger)
			if shopInfoErr != nil {
				logger.WithField("shopInfoErr", shopInfoErr).Error("get shopInfoErr error")
				return &protos.CreateRebateOrderReply{
					Err: &protos.Error{
						Code:    shopInfoErr.Code,
						Message: shopInfoErr.Message,
						Stack:   nil,
					},
				}, nil
			}

			recodeData.FirstType = shopInfo.InviterType
			recodeData.FirstUserId = shopInfo.InviterId
		}
		// 如果是加盟商, 收益归平台
		if firstUserInfo.UserType == rebate.InviterTypeFranchisee {
			recodeData.FirstType = firstUserInfo.InviterType
			recodeData.FirstUserId = firstUserInfo.InviterId
		}
	}

	/**
	 * 1.2 判断如果是加盟商,并且 加盟商是关闭返利的情况则 都归给平台
	 */
	if recodeData.FirstType == rebate.InviterTypeFranchisee {
		// 获取加盟商的消息
		shopInfo, shopInfoErr := GetInviterInfo(recodeData.FirstUserId, req.RebateOrder.ShopId, recodeData.FirstType, logger)
		if shopInfoErr != nil {
			logger.WithField("shopInfoErr", shopInfoErr).Error("get shopInfoErr error")
			return &protos.CreateRebateOrderReply{
				Err: &protos.Error{
					Code:    shopInfoErr.Code,
					Message: shopInfoErr.Message,
					Stack:   nil,
				},
			}, nil
		}

		if !shopInfo.IsRebate {
			// 获取平台数据
			shopInfo, _ := GetInviterInfo(recodeData.FirstUserId, req.RebateOrder.ShopId, rebate.InviterTypeDefault, logger)
			recodeData.FirstType = shopInfo.InviterType
			recodeData.FirstUserId = shopInfo.InviterId
		}

	}

	/*** 二级返利 ***/
	recodeData.SecondType = firstUserInfo.InviterType
	recodeData.SecondUserId = firstUserInfo.InviterId
	recodeData.SecondRebate = recodeData.SecondRebateRate * recodeData.PaidPrice / 100

	// 获取二级返利者信息
	var secondUserInfo *rebate.UserInvited
	secondUserInfo, err = GetInviterInfo(firstUserInfo.InviterId, req.RebateOrder.ShopId, firstUserInfo.InviterType, logger)

	if err != nil {
		logger.WithField("err", err).Error("get firstUserInfo error")
		return &protos.CreateRebateOrderReply{
			Err: &protos.Error{
				Code:    err.Code,
				Message: err.Message,
				Stack:   nil,
			},
		}, nil
	}

	/**
	 * 2.1 判断二级是否开启返利, 如果关闭返利以及是员工/门店, 收益归加盟商
	 */
	logger.
		WithField("IsRebate", firstUserInfo.IsRebate).
		WithField("SecondType", firstUserInfo.UserType).
		WithField("SecondUserId", firstUserInfo.UserId).
		Info("FirstInfo")

	logger.
		WithField("IsRebate", secondUserInfo.IsRebate).
		WithField("SecondType", secondUserInfo.UserType).
		WithField("SecondUserId", secondUserInfo.UserId).
		Info("SecondInfo")

	if !secondUserInfo.IsRebate {
		// 如果是门店, 收益归加盟商
		if secondUserInfo.UserType == rebate.InviterTypeShop {
			recodeData.SecondType = secondUserInfo.InviterType
			recodeData.SecondUserId = secondUserInfo.InviterId
		}
		// 如果是员工, 收益也归加盟商
		if secondUserInfo.UserType == rebate.InviterTypeStaff {
			// 通过员工找上一级的的 Id(门店)找加盟商
			shopInfo, shopInfoErr := GetInviterInfo(secondUserInfo.InviterId, req.RebateOrder.ShopId, secondUserInfo.InviterType, logger)
			if shopInfoErr != nil {
				logger.WithField("shopInfoErr", shopInfoErr).Error("get shopInfoErr error")
				return &protos.CreateRebateOrderReply{
					Err: &protos.Error{
						Code:    shopInfoErr.Code,
						Message: shopInfoErr.Message,
						Stack:   nil,
					},
				}, nil
			}

			recodeData.SecondType = shopInfo.InviterType
			recodeData.SecondUserId = shopInfo.InviterId
		}
		// 如果是加盟商, 收益归平台
		if secondUserInfo.UserType == rebate.InviterTypeFranchisee {
			recodeData.SecondType = secondUserInfo.InviterType
			recodeData.SecondUserId = secondUserInfo.InviterId
		}
	} else {
		/**
		 * 2.2 当一级为加盟商的情况下, 如果一级开启了返利,则二级为平台,如果一级没开启返利二级仍是加盟商
		 */
		if !firstUserInfo.IsRebate {
			recodeData.SecondType = recodeData.FirstType
			recodeData.SecondUserId = recodeData.FirstUserId
		}
	}
	/**
	 * 2.2 判断如果是加盟商,并且 加盟商是关闭返利的情况则 都归给平台
	 */
	if recodeData.SecondType == rebate.InviterTypeFranchisee {
		// 获取加盟商的消息
		shopInfo, shopInfoErr := GetInviterInfo(recodeData.SecondUserId, req.RebateOrder.ShopId, recodeData.SecondType, logger)
		if shopInfoErr != nil {
			logger.WithField("shopInfoErr", shopInfoErr).Error("get shopInfoErr error")
			return &protos.CreateRebateOrderReply{
				Err: &protos.Error{
					Code:    shopInfoErr.Code,
					Message: shopInfoErr.Message,
					Stack:   nil,
				},
			}, nil
		}

		if !shopInfo.IsRebate {
			// 获取平台数据
			shopInfo, _ := GetInviterInfo(recodeData.FirstUserId, req.RebateOrder.ShopId, rebate.InviterTypeDefault, logger)
			recodeData.SecondType = shopInfo.InviterType
			recodeData.SecondUserId = shopInfo.InviterId
		}
	}

	returnData, returnDataErr := CreateOrderRecode(recodeData)
	if returnDataErr != nil {
		logger.WithError(returnDataErr).Error("create order rebate recode error")
		return &protos.CreateRebateOrderReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeCreateRebateOrderListError,
				Message: constant.ErrorMessageCreateOrderListError,
				Stack:   nil,
			},
		}, nil
	}
	return &protos.CreateRebateOrderReply{
		RebateOrder: returnData,
	}, nil
}
