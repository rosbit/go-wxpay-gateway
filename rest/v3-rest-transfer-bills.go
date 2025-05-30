package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3-selfdev"
	"net/http"
)

// POST /v3/transfer-bills
// POST Body:
// {
//  "payApp": "name-of-app-in-wxpay-gateway",
//  "appid": "APP、小程序、公众号、企业号corpid即为此AppID",
//  "out_bills_no": "商户系统内部的商家单号，只能由数字、大小写字母组成，在商户系统内部唯一",
//  "transfer_scene_id": " 该笔转账使用的转账场景，可前往“商户平台-产品中心-商家转账”中申请。如：1000（现金营销），1006（企业报销",
//  "openid": "用户在商户appid下的唯一标识",
//  "user_name": "收款方真实姓名，转账金额 >= 2,000元时，该笔明细必须填写",
//  "transfer_amount": 1024, // 转账金额单位为“分”
//  "transfer_remark": "转账备注，用户收款时可见该备注，UTF8编码，最多允许32个字符",
//  "notify_url": "异步接收微信支付结果通知的回调地址，必须为HTTPS，不能携带参数"
//  "user_recv_perception": "用户收款时感知到的收款原因将根据转账场景自动展示默认内容。如有其他展示需求，可在本字段传入。",
//  "transfer_scene_report_infos": [ // 各转账场景下需报备的内容，商户需要按照所属转账场景规则传参
//    {
//       "info_type": "不能超过15个字符，商户所属转账场景下的信息类型，此字段内容为固定值，需严格按照转账场景报备信息字段说明传参",
//       "info_content": "不能超过32个字符，商户所属转账场景下的信息内容，商户可按实际业务场景自定义传参，需严格按照转账场景报备信息字段说明传参"
//    },
//    { ...}
//  ]
// }
func V3TransferBills(c *mgin.Context) {
	var params struct {
		PayApp   string `json:"payApp"`
		RealUserName string `json:"user_name"`
		v3sd.CreateTransferBillsRequest
	}

	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	params.CreateTransferBillsRequest.UserName = params.RealUserName

	resp, err := v3sd.CreateTransferBills(params.PayApp, &params.CreateTransferBillsRequest)
	if err != nil {
		sendResultWithMsg(c, true, nil, nil, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": resp,
	})
}

