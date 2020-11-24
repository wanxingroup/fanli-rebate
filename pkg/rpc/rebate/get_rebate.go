package rebate

import (
	"fmt"

	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func (_ Controller) GetRebate(ctx context.Context, req *protos.GetRebateRequest) (*protos.GetRebateReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateGetRebate(req)

	if err != nil {
		return &protos.GetRebateReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	// 获取原来的数据
	rebate, rebateErr := GetRebateByIteam(req.ItemId, uint8(req.ItemType))
	if rebateErr != nil {

		if rebateErr == gorm.ErrRecordNotFound {

			logger.WithError(err).Error("get recode is empty")
			return &protos.GetRebateReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeGetRebateEmpty,
					Message: constant.ErrorMessageGetRebateEmpty,
					Stack:   nil,
				},
			}, nil

		} else {
			logger.WithError(err).Error("get rebate error")
			return &protos.GetRebateReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeGetRebateError,
					Message: constant.ErrorMessageGetRebateError,
					Stack:   nil,
				},
			}, nil
		}

	}
	return &protos.GetRebateReply{
		Rebate: ConvertModelToProtobuf(rebate),
	}, nil
}

func validateGetRebate(req *protos.GetRebateRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ItemId, ItemIdRule...),
		validation.Field(&req.ItemType, ItemTypeRule...),
	)
}
