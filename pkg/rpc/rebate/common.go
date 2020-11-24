package rebate

import (
	"context"
	"strconv"
	"time"

	cardProtos "dev-gitlab.wanxingrowth.com/fanli/card/pkg/rpc/protos"
	merchantProtos "dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	userProtos "dev-gitlab.wanxingrowth.com/fanli/user/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"

	baseDb "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/model/rebate"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

const (
	DefaultPage        = 1
	DefaultPageSize    = 20
	DefaultMaxPageSize = 100
)

func GetRebateByIteam(itemId uint64, itemType uint8) (*rebate.RebateItem, error) {
	var record rebate.RebateItem
	err := database.GetDB(constant.DatabaseConfigKey).Where("itemId = ? and itemType = ?", itemId, itemType).First(&record).Error

	if err != nil {
		log.GetLogger().WithField("rebate", record).WithError(err).Error("get record error")
		return nil, err
	}

	return &record, nil
}

// 写入返利
func CreateRebate(request *protos.Rebate) (*protos.Rebate, error) {

	record := &rebate.RebateItem{
		ItemId:            request.ItemId,
		ItemType:          uint8(request.ItemType),
		FirstRebateRatio:  request.FirstRebateRatio,
		SecondRebateRatio: request.SecondRebateRatio,
		CreatedBy:         request.CreatedBy,
		UpdatedBy:         request.CreatedBy,
	}
	err := database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {
		log.GetLogger().WithField("rebate", record).WithError(err).Error("create record error")
		return nil, err
	}

	return request, nil
}

// 更新返利
func ModifyRebate(request *protos.Rebate) (*protos.Rebate, error) {

	record := &rebate.RebateItem{
		FirstRebateRatio:  request.FirstRebateRatio,
		SecondRebateRatio: request.SecondRebateRatio,
		UpdatedBy:         request.UpdatedBy,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Model(rebate.RebateItem{}).Where("itemId = ? and itemType = ?", request.ItemId, request.ItemType).Update(record).Error
	if err != nil {
		log.GetLogger().WithField("rebase", record).WithError(err).Error("modify record error")
		return nil, err
	}

	return request, nil
}

func ConvertErrorToProtobuf(err error) *protos.Error {

	if validationError, ok := err.(validation.Error); ok {
		errorCode, convertError := strconv.Atoi(validationError.Code())
		if convertError != nil {
			errorCode = errors.CodeServerInternalError
		}
		return &protos.Error{
			Code:    int64(errorCode),
			Message: validationError.Error(),
		}
	}

	return &protos.Error{
		Code:    errors.CodeServerInternalError,
		Message: err.Error(),
	}
}

func ConvertModelToProtobuf(data *rebate.RebateItem) *protos.Rebate {
	return &protos.Rebate{
		ItemId:            data.ItemId,
		ItemType:          uint32(data.ItemType),
		FirstRebateRatio:  data.FirstRebateRatio,
		SecondRebateRatio: data.SecondRebateRatio,
	}
}

// 获取订单统计数据
func GetUserOrderListByCondition(conditions map[string]interface{}, pageData map[string]uint64) ([]*rebate.RebateOrder, uint64, error) {
	db := database.GetDB(constant.DatabaseConfigKey).Model(&rebate.RebateOrder{})
	if shopId, has := conditions["shopId"]; has {
		db = db.Where("shopId = ?", shopId)
	}

	if orderId, has := conditions["orderId"]; has {
		db = db.Where("orderId = ?", orderId)
	}

	if userType, has := conditions["userType"]; has {
		db = db.Where("firstType = ? or secondType = ?", userType, userType)
	}

	if userId, has := conditions["userId"]; has {
		db = db.Where("firstUserId = ? or secondUserId = ?", userId, userId)
	}

	if paidTimeStart, has := conditions["paidTimeStart"]; has {
		db = db.Where("paidTime > ?", paidTimeStart)
	}

	if paidTimeEnd, has := conditions["paidTimeEnd"]; has {
		db = db.Where("paidTime < ?", paidTimeEnd)
	}

	db = db.Order("orderId desc")

	results := make([]*rebate.RebateOrder, 0, pageData["pageSize"])
	var count uint64
	err := baseDb.FindPage(db, pageData, &results, &count)
	return results, count, err
}

func GetUserOrderListTotalByCondition(conditions map[string]interface{}, statisticsResult *rebate.StatisticsOrderResult) {
	db := database.GetDB(constant.DatabaseConfigKey).Model(&rebate.RebateOrder{}).
		Select("sum(paidPrice) as payMoneyTotal, sum(firstRebate) as firstRebateTotal, sum(secondRebate) as secondRebateTotal")

	if shopId, has := conditions["shopId"]; has {
		db = db.Where("shopId = ?", shopId)
	}

	if orderId, has := conditions["orderId"]; has {
		db = db.Where("orderId = ?", orderId)
	}

	if userType, has := conditions["userType"]; has {
		db = db.Where("firstType = ? or secondType = ?", userType, userType)
	}

	if userId, has := conditions["userId"]; has {
		db = db.Where("firstUserId = ? or secondUserId = ?", userId, userId)
	}

	if paidTimeStart, has := conditions["paidTimeStart"]; has {
		db = db.Where("paidTime > ?", paidTimeStart)
	}

	if paidTimeEnd, has := conditions["paidTimeEnd"]; has {
		db = db.Where("paidTime < ?", paidTimeEnd)
	}

	err := db.Scan(statisticsResult).Error

	if err != nil {
		log.GetLogger().WithError(err).Error("get user list statistics error")
	}
}

// 创建记录
func CreateOrderRecode(request *protos.RebateOrder) (*protos.RebateOrder, error) {

	paidTime, err := time.ParseInLocation("2006-01-02 15:04:05", request.PaidTime, time.Local)
	if err != nil {
		log.GetLogger().WithError(err).Error("convert time error")
		paidTime = time.Now()
	}

	record := &rebate.RebateOrder{
		OrderId:          request.OrderId,
		ShopId:           request.ShopId,
		UserId:           request.UserId,
		Mobile:           request.Mobile,
		PaidPrice:        request.PaidPrice,
		FirstRebate:      request.FirstRebate,
		FirstUserId:      request.FirstUserId,
		FirstType:        uint8(request.FirstType),
		FirstRebateRate:  request.FirstRebateRate,
		SecondRebateRate: request.SecondRebateRate,
		SecondRebate:     request.SecondRebate,
		SecondUserId:     request.SecondUserId,
		SecondType:       uint8(request.SecondType),
		PaidTime:         paidTime,
		ItemType:         uint8(request.ItemType),
		ItemId:           request.ItemId,
	}

	err = database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {
		log.GetLogger().WithField("rebate", record).WithError(err).Error("create order list record error")
		return nil, err
	}

	return request, nil
}

// 获取邀请者信息的信息, 以及上一级的用户类型以及Id
func GetInviterInfo(userId, shopId uint64, userType uint32, logger *logrus.Entry) (rebData *rebate.UserInvited, proErr *protos.Error) {
	returnData := &rebate.UserInvited{
		UserMobile:  "",
		UserId:      userId,
		InviterId:   0,
		InviterType: 0,
		IsRebate:    false,
		UserName:    "",
	}
	switch userType {
	case rebate.InviterTypeDefault:
		logger.Info("InviterTypeDefault")

		returnData = &rebate.UserInvited{
			UserId:      0,
			InviterId:   0,
			InviterType: rebate.InviterTypeDefault,
			IsRebate:    true,
			UserType:    rebate.InviterTypeDefault,
			UserName:    "平台",
		}
	case rebate.InviterTypeFranchisee:
		logger.Info("InviterTypeFranchisee")

		api := client.GetMerchantService()
		if api == nil {
			log.GetLogger().Error("get merchant rpc service client is nil")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeInternalServiceConnectionError,
				Message: constant.ErrorMessageInternalServiceConnectionError,
			}
		}
		reply, err := api.GetMerchant(context.Background(), &merchantProtos.GetMerchantRequest{
			MerchantId: userId, // 库里的所有邀请者字段都为 userId
		})

		if err != nil {
			logger.WithError(err).Error("call merchant service get list error")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeCallMerchantServiceError,
				Message: constant.ErrorMessageCallServiceMerchantError,
			}
		}

		if reply.Err != nil {
			logger.WithField("replyError", *reply.Err).Info("merchant service reply error")

			return nil, &protos.Error{
				Code:    reply.Err.GetCode(),
				Message: reply.Err.GetMessage(),
			}
		}

		returnData = &rebate.UserInvited{
			UserId:      reply.Merchant.MerchantId,
			InviterId:   0,
			InviterType: rebate.InviterTypeDefault,
			IsRebate:    reply.Merchant.IsRebate,
			UserType:    rebate.InviterTypeFranchisee,
			UserName:    reply.Merchant.Name,
		}
	case rebate.InviterTypeShop:
		logger.Info("InviterTypeShop")

		api := client.GetShopService()
		if api == nil {
			log.GetLogger().Error("get shop rpc service client is nil")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeInternalServiceConnectionError,
				Message: constant.ErrorMessageInternalServiceConnectionError,
			}
		}
		reply, err := api.GetShopInfo(context.Background(), &merchantProtos.GetShopInfoRequest{
			ShopId: userId, // 库里的所有邀请者字段都为 userId
		})

		if err != nil {
			logger.WithError(err).Error("call shop service get staff error")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeCallMerchantServiceError,
				Message: constant.ErrorMessageCallServiceMerchantError,
			}
		}

		if reply.Err != nil {
			logger.WithField("replyError", *reply.Err).Info("shop service reply error")

			return nil, &protos.Error{
				Code:    reply.Err.GetCode(),
				Message: reply.Err.GetMessage(),
			}

		}

		returnData = &rebate.UserInvited{
			InviterId:   reply.ShopInfo.MerchantId,
			UserId:      reply.ShopInfo.ShopId,
			InviterType: rebate.InviterTypeFranchisee,
			IsRebate:    reply.ShopInfo.IsRebate,
			UserType:    rebate.InviterTypeShop,
			UserName:    reply.ShopInfo.Name,
		}
	case rebate.InviterTypeStaff:
		logger.Info("InviterTypeStaff")
		api := client.GetStaffService()
		if api == nil {
			log.GetLogger().Error("get staff rpc service client is nil")

			return nil, &protos.Error{
				Code:    constant.ErrorCodeInternalServiceConnectionError,
				Message: constant.ErrorMessageInternalServiceConnectionError,
			}

		}
		reply, err := api.Get(context.Background(), &merchantProtos.GetStaffInfoRequest{
			StaffId: userId,
		})

		if err != nil {
			logger.WithError(err).Error("call staff service get staff error")

			return nil, &protos.Error{
				Code:    constant.ErrorCodeCallMerchantServiceError,
				Message: constant.ErrorMessageCallServiceMerchantError,
			}
		}

		if reply.Err != nil {
			logger.WithField("replyError", *reply.Err).Info("staff service reply error")
			return nil, &protos.Error{
				Code:    reply.Err.GetCode(),
				Message: reply.Err.GetMessage(),
			}
		}

		returnData = &rebate.UserInvited{
			InviterId:   reply.StaffInfo.ShopId,
			UserId:      reply.StaffInfo.StaffId,
			InviterType: rebate.InviterTypeShop,
			IsRebate:    reply.StaffInfo.IsRebate,
			UserType:    rebate.InviterTypeStaff,
			UserName:    reply.StaffInfo.Name,
		}
	case rebate.InviterTypeUser:
		logger.Info("InviterTypeUser")
		api := client.GetUserService()
		if api == nil {
			logger.Error("get user rpc service client is nil")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeInternalServiceConnectionError,
				Message: constant.ErrorMessageInternalServiceConnectionError,
			}

		}

		reply, err := api.GetUserInfo(context.Background(), &userProtos.GetUserInfoRequest{
			ShopId: shopId,
			UserId: userId,
		})

		log.GetLogger().WithField("reply", reply).Info("get user reply")

		if err != nil {
			logger.WithError(err).Error("call oss service error")
			return nil, &protos.Error{
				Code:    constant.ErrorCodeCallMerchantServiceError,
				Message: constant.ErrorMessageCallServiceMerchantError,
			}

		}

		if reply.Err != nil {

			logger.WithField("replyError", *reply.Err).Info("oss service reply error")

			return nil, &protos.Error{
				Code:    reply.Err.GetCode(),
				Message: reply.Err.GetMessage(),
			}
		}

		returnData = &rebate.UserInvited{
			UserMobile:  reply.UserInfo.Mobile,
			InviterId:   reply.UserInfo.InviterId,
			UserId:      reply.UserInfo.UserId,
			InviterType: reply.UserInfo.InviterType,
			IsRebate:    true,
			UserType:    rebate.InviterTypeUser,
			UserName:    reply.UserInfo.Name,
		}
	}
	return returnData, nil
}

// 获取当前用户是否购买过权益卡
func GetUserIsHaveCard(logger *logrus.Entry, userId uint64) bool {
	api := client.GetCardService()
	if api == nil {
		logger.Error("get user rpc service client is nil")
		return false
	}

	reply, err := api.GetUserCardList(context.Background(), &cardProtos.GetUserCardListRequest{
		UserId: userId,
	})

	log.GetLogger().WithField("reply", reply).Info("get user reply")

	if err != nil {
		logger.WithError(err).Error("call oss service error")
		return false
	}

	if reply.Err != nil {

		logger.WithField("replyError", *reply.Err).Info("oss service reply error")
		return false
	}

	if len(reply.UserCardList) > 0 {
		return true
	}

	return false
}
