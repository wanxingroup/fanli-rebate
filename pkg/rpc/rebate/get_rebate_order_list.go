package rebate

import (
	"fmt"

	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/model/rebate"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func (_ Controller) GetRebateOrderList(ctx context.Context, req *protos.GetRebateOrderListRequest) (*protos.GetRebateOrderListReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	var page = uint64(DefaultPage)
	var pageSize = uint64(DefaultPageSize)

	if req.GetPage() > 0 {
		page = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
	}

	if !validation.IsEmpty(req.GetOrderId()) {
		conditions["orderId"] = req.GetOrderId()
	}

	if !validation.IsEmpty(req.GetUserType()) {
		conditions["userType"] = req.GetUserType()
	}

	if !validation.IsEmpty(req.GetUserId()) {
		conditions["userId"] = req.GetUserId()
	}

	if !validation.IsEmpty(req.GetPaidTimeStart()) && !validation.IsEmpty(req.GetPaidTimeEnd()) {
		conditions["paidTimeStart"] = req.GetPaidTimeStart()
		conditions["paidTimeEnd"] = req.GetPaidTimeEnd()
	}

	pageData := map[string]uint64{
		"page":     page,
		"pageSize": pageSize,
	}

	list, count, err := GetUserOrderListByCondition(conditions, pageData)

	if err != nil {
		logger.WithError(err).Error("get list distribution user error")
		return &protos.GetRebateOrderListReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetRebateOrderListError,
				Message: constant.ErrorMessageGetOrderListError,
				Stack:   nil,
			},
		}, nil
	}
	// 循环获取用户/平台名字

	returnList := make([]*protos.RebateOrder, len(list))
	for key, data := range list {
		rebateData := &protos.RebateOrder{
			OrderId:      data.OrderId,
			UserId:       data.UserId,
			PaidPrice:    data.PaidPrice,
			FirstRebate:  data.FirstRebate,
			FirstUserId:  data.FirstUserId,
			FirstType:    uint32(data.FirstType),
			SecondRebate: data.SecondRebate,
			SecondUserId: data.SecondUserId,
			SecondType:   uint32(data.SecondType),
			ShopId:       data.ShopId,
			PaidTime:     data.PaidTime.Format("2006-01-02 15:04:05"),
			ItemType:     uint32(data.ItemType),
			ItemId:       data.ItemId,
		}

		// 获取用户昵称
		userInfo, err := GetInviterInfo(data.UserId, data.ShopId, rebate.InviterTypeUser, logger)
		if err != nil {
			logger.WithField("err", err).WithField("FirstUserId", rebateData.FirstUserId).WithField("FirstType", rebateData.FirstType).Error("get first UserInfo error")
			rebateData.UserName = "此人-查询错误"
		} else {
			rebateData.UserName = userInfo.UserName
		}

		// 获取一级返利昵称
		firstUserInfo, err := GetInviterInfo(rebateData.FirstUserId, rebateData.ShopId, rebateData.FirstType, logger)
		if err != nil {
			logger.WithField("err", err).WithField("FirstUserId", rebateData.FirstUserId).WithField("FirstType", rebateData.FirstType).Error("get first UserInfo error")
			rebateData.FirstUserName = "此人-查询错误"
		} else {
			rebateData.FirstUserName = firstUserInfo.UserName
		}

		// 获取二级返利昵称
		secondUserInfo, err := GetInviterInfo(rebateData.SecondUserId, rebateData.ShopId, rebateData.SecondType, logger)
		if err != nil {
			logger.WithField("err", err).WithField("SecondUserId", rebateData.SecondUserId).WithField("SecondType", rebateData.SecondType).Error("get second UserInfo error")
			rebateData.SecondUserName = "此人-查询错误"
		} else {
			rebateData.SecondUserName = secondUserInfo.UserName
		}

		returnList[key] = rebateData
	}

	// 获取统计数据
	var statisticsResult rebate.StatisticsOrderResult
	GetUserOrderListTotalByCondition(conditions, &statisticsResult)

	return &protos.GetRebateOrderListReply{
		RebateTotal: statisticsResult.FirstRebateTotal + statisticsResult.SecondRebateTotal,
		RebateOrder: returnList,
		PaidTotal:   statisticsResult.PayMoneyTotal,
		Count:       count,
	}, nil

}
