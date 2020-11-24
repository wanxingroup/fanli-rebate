package rebate

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/utils/log"
)

func (_ Controller) BatchRebate(ctx context.Context, req *protos.BatchRebateRequest) (*protos.BatchRebateReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	// 批量更新

	for _, rebate := range req.Rebate {
		// 获取原来的数据
		_, rebateErr := GetRebateByIteam(rebate.ItemId, uint8(rebate.ItemType))

		if rebateErr != nil {
			if rebateErr == gorm.ErrRecordNotFound {
				// 写入
				_, err := CreateRebate(rebate)
				if err != nil {
					logger.WithError(err).Error("create rebate error")
					return &protos.BatchRebateReply{
						Err: &protos.Error{
							Code:    constant.ErrorCodeCreateRebateError,
							Message: constant.ErrorMessageCreateRebateError,
							Stack:   nil,
						},
					}, nil
				}

			} else {
				logger.Error("get rebate error")
				return &protos.BatchRebateReply{
					Err: &protos.Error{
						Code:    constant.ErrorCodeGetRebateError,
						Message: constant.ErrorMessageGetRebateError,
						Stack:   nil,
					},
				}, nil
			}
		} else {
			// 更新 返利
			_, err := ModifyRebate(rebate)
			if err != nil {
				logger.WithError(err).Error("modify distribution error")
				return &protos.BatchRebateReply{
					Err: &protos.Error{
						Code:    constant.ErrorCodeModifyRebateError,
						Message: constant.ErrorMessageModifyRebateError,
						Stack:   nil,
					},
				}, nil
			}
		}
	}

	// 返回 结果
	return &protos.BatchRebateReply{
		Err: nil,
	}, nil

}
