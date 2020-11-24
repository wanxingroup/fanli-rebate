package rebate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/rebate/pkg/constant"
)

var ItemIdRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeErrorCodeItemIdRequired, constant.ErrorMessageErrorCodeItemIdRequired)),
}

var ItemTypeRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeErrorCodeItemTypeRequired, constant.ErrorMessageErrorCodeItemTypeRequired)),
}

var FirstRebateRatioRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeErrorCodeFirstRebateRatioRuleRequired, constant.ErrorMessageErrorCodeFirstRebateRatioRuleRequired)),
}

var SecondRebateRatioRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeErrorCodeSecondRebateRatioRuleRequired, constant.ErrorMessageErrorCodeSecondRebateRatioRuleRequired)),
}
