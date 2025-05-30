// 商家转账相关数据定义
// 参考文档: 
//  1. 发起转账: https://pay.weixin.qq.com/doc/v3/merchant/4012716434
package v3sd

import (
	"time"
)

// 各转账场景下需报备的内容，商户需要按照所属转账场景规则传参
type TransferSceneReportInfo struct {
	InfoType string `json:"info_type"` // 不能超过15个字符，商户所属转账场景下的信息类型，此字段内容为固定值，需严格按照转账场景报备信息字段说明传参
	InfoContent string `json:"info_content"`// 不能超过32个字符，商户所属转账场景下的信息内容，商户可按实际业务场景自定义传参，需严格按照转账场景报备信息字段说明传参。
}

type CreateTransferBillsRequest struct {
	Appid string `json:"appid"` // 是微信开放平台和微信公众平台为开发者的应用程序(APP、小程序、公众号、企业号corpid即为此AppID)提供的一个唯一标识。此处，可以填写这四种类型中的任意一种APPID，但请确保该appid与商户号有绑定关系。
	OutBillNo string `json:"out_bill_no"` // 商户系统内部的商家单号，要求此参数只能由数字、大小写字母组成，在商户系统内部唯一
	TransferSceneId string `json:"transfer_scene_id"` // 该笔转账使用的转账场景，可前往“商户平台-产品中心-商家转账”中申请。如：1000（现金营销），1006（企业报销）等
	Openid string `json:"openid"` // 用户在商户appid下的唯一标识。发起转账前需获取到用户的OpenID
	UserName string `json:"-"` // 收款方真实姓名。需要加密传入，支持标准RSA算法和国密算法，公钥由微信侧提供。
                                       // 转账金额 >= 2,000元时，该笔明细必须填写
									   // 若商户传入收款用户姓名，微信支付会校验收款用户与输入姓名是否一致，并提供电子回单
	TransferAmount int32 `json:"transfer_amount"` // 转账金额单位为“分”。
	TransferRemark string `json:"transfer_remark"` // 转账备注，用户收款时可见该备注信息，UTF8编码，最多允许32个字符
	NotifyUrl string `json:"notify_url,omitempty"` // 异步接收微信支付结果通知的回调地址，通知url必须为公网可访问的URL，必须为HTTPS，不能携带参数。
	UserRecvPerception string `json:"user_recv_perception,omitempty"` //  用户收款时感知到的收款原因将根据转账场景自动展示默认内容。如有其他展示需求，可在本字段传入。各场景展示的默认内容和支持传入的内容，可查看产品文档了解
	TransferSceneReportInfos []TransferSceneReportInfo `json:"transfer_scene_report_infos"`// 各转账场景下需报备的内容，商户需要按照所属转账场景规则传参
}

/*
// 发起转账返回值
type CreateTransferBillsResponse struct {
	OutBillNo string `json:"out_bill_no"` // 商户系统内部的"商家单号"，要求此参数只能由数字、大小写字母组成，在商户系统内部唯一
	TransferBillNo string `json:"transfer_bill_no"` // 微信转账单号，微信商家转账系统返回的唯一标识
	CreateTime string `json:"create_time"` // 单据创建时间, 单据受理成功时返回，按照使用rfc3339所定义的格式，格式为yyyy-MM-DDThh:mm:ss+TIMEZONE
	State string `json:"state"` // 商家转账订单状态
	                            // ACCEPTED: 转账已受理
								// PROCESSING: 转账锁定资金中。如果一直停留在该状态，建议检查账户余额是否足够，如余额不足，可充值后再原单重试。
								// WAIT_USER_CONFIRM: 待收款用户确认，可拉起微信收款确认页面进行收款确认
								// TRANSFERING: 转账中，可拉起微信收款确认页面再次重试确认收款
								// SUCCESS: 转账成功
								// FAIL: 转账失败
								// CANCELING: 商户撤销请求受理成功，该笔转账正在撤销中
								// CANCELLED: 转账撤销完成
	PackageInfo string `json:"package_info"` // 跳转微信支付收款页的package信息，APP调起用户确认收款或者JSAPI调起用户确认收款 时需要使用的参数。

}

// 撤销转账返回值
type CancelTransferBillsResponse struct {
	OutBillNo string `json:"out_bill_no"` // 商户系统内部的"商家单号"，要求此参数只能由数字、大小写字母组成，在商户系统内部唯一
	TransferBillNo string `json:"transfer_bill_no"` // 微信转账单号，微信商家转账系统返回的唯一标识
	UpdateTime string `json:"update_time"` // 最后一次单据状态变更时间，按照使用rfc3339所定义的格式，格式为yyyy-MM-DDThh:mm:ss+TIMEZONE
	State string `json:"state"` // 商家转账订单状态
	                            // ACCEPTED: 转账已受理
								// PROCESSING: 转账锁定资金中。如果一直停留在该状态，建议检查账户余额是否足够，如余额不足，可充值后再原单重试。
								// WAIT_USER_CONFIRM: 待收款用户确认，可拉起微信收款确认页面进行收款确认
								// TRANSFERING: 转账中，可拉起微信收款确认页面再次重试确认收款
								// SUCCESS: 转账成功
								// FAIL: 转账失败
								// CANCELING: 商户撤销请求受理成功，该笔转账正在撤销中
								// CANCELLED: 转账撤销完成
}*/

type VerifyBody struct {
	Id string `json:"id"`
	CreateTime time.Time `json:"create_time"`
	ResourceType string `json:"resource_type"`
	EventType string `json:"event_type"`
	Summary string `json:"summary"`
    Resource struct {
		OriginalType string `json:"original_type"`
		Algorithm string `json:"algorithm"`
		Ciphertext string `json:"ciphertext"`
		AssociatedData string `json:"associated_data"`
		Nonce string `json:"nonce"`
	} `json:"resource"`
}

/*
// 解密后的单据
type Bill struct {
	OutBillNo string `json:"out_bill_no"` // 【商户单号】商户系统内部的商家单号，在商户系统内部唯一
	TransferBillNo string `json:"transfer_bill_no"` // 商家转账订单号】微信单号，微信商家转账系统返回的唯一标识
	State string `json:"state"` // 【单据状态】商家转账订单状态
                                // ACCEPTED：单据已受理
								// PROCESSING：单据处理中，转账结果尚未明确，如一直处于此状态，建议检查账户余额是否足够
								// WAIT_USER_CONFIRM：待收款用户确认，可拉起微信收款确认页面进行收款确认
								// TRANSFERING：转账中，转账结果尚未明确，可拉起微信收款确认页面再次重试确认收款
								// SUCCESS： 转账成功
								// FAIL： 转账失败
								// CANCELING： 撤销中
								// CANCELLED： 已撤销
	MchId string `json:"mch_id"` // 【商户号】微信支付分配的商户号
	TransferAmount int32 `json:"transfer_amount"` // 【转账金额】转账总金额，单位为“分”
	OpenId string `json:"openid"` // 【收款用户OpenID】用户在商户appid下的唯一标识
	FailReason string `json:"fail_reason:"` // 【失败原因】单已失败或者已退资金时，会返回订单失败原因
	CreateTime time.Time `json:"create_time"` // 【单据创建时间】遵循rfc3339标准格式，格式为yyyy-MM-DDTHH:mm:ss+TIMEZONE，yyyy-MM-DD表示年月日，T出现在字符串中，表示time元素的开头，HH:mm:ss.表示时分秒，TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。例如：2015-05-20T13:29:35+08:00表示北京时间2015年05月20日13点29分35秒。
	UpdateTime time.Time `json:"update_time"` // 【最后一次状态变更时间】格式同CreateTime
}
*/
