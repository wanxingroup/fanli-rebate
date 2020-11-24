package constant

const (
	ErrorCodeErrorCodeItemIdRequired                   = "419001"
	ErrorMessageErrorCodeItemIdRequired                = "itemId 为必填"
	ErrorCodeErrorCodeItemTypeRequired                 = "419002"
	ErrorMessageErrorCodeItemTypeRequired              = "类型 为必填"
	ErrorCodeErrorCodeFirstRebateRatioRuleRequired     = "419003"
	ErrorMessageErrorCodeFirstRebateRatioRuleRequired  = "一级返利 必填"
	ErrorCodeErrorCodeSecondRebateRatioRuleRequired    = "419004"
	ErrorMessageErrorCodeSecondRebateRatioRuleRequired = "二级返利 必填"
)
const (
	ErrorCodeTransactionError                  = 519005
	ErrorMessageTransactionError               = "数据库事务出错"
	ErrorCodeGetRebateError                    = 519006
	ErrorMessageGetRebateError                 = "获取返利活动失败"
	ErrorCodeCreateRebateError                 = 519007
	ErrorMessageCreateRebateError              = "创建返利活动失败"
	ErrorCodeModifyRebateError                 = 519008
	ErrorMessageModifyRebateError              = "更新返利活动失败"
	ErrorCodeGetRebateEmpty                    = 519009
	ErrorMessageGetRebateEmpty                 = "返利数据为空"
	ErrorCodeGetRebateOrderListError           = 519010
	ErrorMessageGetOrderListError              = "返利订单获取失败"
	ErrorCodeCreateRebateOrderListError        = 519011
	ErrorMessageCreateOrderListError           = "创建返利订单失败"
	ErrorCodeInternalServiceConnectionError    = 519012
	ErrorMessageInternalServiceConnectionError = "内部服务连接出错"
	ErrorCodeCallMerchantServiceError          = 519013
	ErrorMessageCallServiceMerchantError       = "调用商家服务出错"
)
