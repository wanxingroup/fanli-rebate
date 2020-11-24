package rebate

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func (_ Controller) SetRebate(ctx context.Context, req *protos.SetRebateRequest) (*protos.SetRebateReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateSetRebate(req)

	if err != nil {
		return &protos.SetRebateReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	// 获取原来的数据
	_, rebateErr := GetRebateByIteam(req.Rebate.ItemId, uint8(req.Rebate.ItemType))

	if rebateErr != nil {
		if rebateErr == gorm.ErrRecordNotFound {
			// 写入
			rebate, err := CreateRebate(req.Rebate)
			if err != nil {
				logger.WithError(err).Error("create rebate error")
				return &protos.SetRebateReply{
					Err: &protos.Error{
						Code:    constant.ErrorCodeCreateRebateError,
						Message: constant.ErrorMessageCreateRebateError,
						Stack:   nil,
					},
				}, nil
			}

			// 返回 结果
			return &protos.SetRebateReply{
				Rebate: rebate,
			}, nil

		} else {
			logger.WithError(err).Error("get rebate error")
			return &protos.SetRebateReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeGetRebateError,
					Message: constant.ErrorMessageGetRebateError,
					Stack:   nil,
				},
			}, nil
		}
	}

	// 更新 返利
	rebate, err := ModifyRebate(req.Rebate)
	if err != nil {
		logger.WithError(err).Error("modify rebate error")
		return &protos.SetRebateReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeModifyRebateError,
				Message: constant.ErrorMessageModifyRebateError,
				Stack:   nil,
			},
		}, nil
	}

	// 返回 结果
	return &protos.SetRebateReply{
		Rebate: rebate,
	}, nil

}

func validateSetRebate(req *protos.SetRebateRequest) error {
	return validation.ValidateStruct(req.Rebate,
		validation.Field(&req.Rebate.ItemId, ItemIdRule...),
		validation.Field(&req.Rebate.ItemType, ItemTypeRule...),
	)
}
