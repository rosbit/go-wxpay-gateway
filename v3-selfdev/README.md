## 商家转账API接口

### 产品介绍

- https://pay.weixin.qq.com/doc/v3/merchant/4012711988
- 商家转账升级版本已于2025年1月15日正式上线。新版本无收款用户管理、商户出资确认功能。

### 接口服务

#### 1. 发起商家转账

- 商家转账用户确认模式下，用户申请收款时，商户可通过此接口申请创建转账单

- 参考文档: https://pay.weixin.qq.com/doc/v3/merchant/4012716434

- 请求接口: POST /v3/transfer-bills

- Body:
  
  ```json
  {
     "payApp": "微信支付网关的应用名，如jiaomu-net",
     "appid": "APP、小程序、公众号、企业号corpid即为此AppID",
     "out_bills_no": "商户系统内部的商家单号，只能由数字、大小写字母组成，在商户系统内部唯一",
     "transfer_scene_id": " 该笔转账使用的转账场景，可前往“商户平台-产品中心-商家转账”中申请。如：1000（现金营销），1006（企业报销",
     "openid": "用户在商户appid下的唯一标识",
     "user_name": "收款方真实姓名，转账金额 >= 2,000元时，该笔明细必须填写",
     "transfer_amount": 1024, // 转账金额单位为“分”
     "transfer_remark": "转账备注，用户收款时可见该备注，UTF8编码，最多允许32个字符",
     "notify_url": "异步接收微信支付结果通知的回调地址，必须为HTTPS，不能携带参数"
     "user_recv_perception": "用户收款时感知到的收款原因将根据转账场景自动展示默认内容。如有其他展示需求，可在本字段传入。",
     "transfer_scene_report_infos": [ // 各转账场景下需报备的内容，商户需要按照所属转账场景规则传参
        {
          "info_type": "不能超过15个字符，商户所属转账场景下的信息类型，此字段内容为固定值，需严格按照转账场景报备信息字段说明传参",
          "info_content": "不能超过32个字符，商户所属转账场景下的信息内容，商户可按实际业务场景自定义传参，需严格按照转账场景报备信息字段说明传参"
        },
        { ...}
      ]
  }
  ```

- 成功返回
  
  ```json
  {
     "code": 200,
     "msg": "OK",
     "result: {
        "out_bill_no" : "plfk2020042013",
        "transfer_bill_no" : "1330000071100999991182020050700019480001",
        "create_time" : "2015-05-20T13:29:35.120+08:00",
        "state" : "ACCEPTED"
        "package_info" : "affffddafdfafddffda=="
     },
  }
  ```

#### 2. 撤销商家转账

- 商户通过转账接口发起付款后，在用户确认收款之前可以通过该接口撤销付款。该接口返回成功仅表示撤销请求已受理，系统会异步处理退款等操作，以最终查询单据返回状态为准。

- 参考文档: https://pay.weixin.qq.com/doc/v3/merchant/4012716458

- 请求接口: POST /v3/cancel-transfer-bills

- Body:
  
  ```json
  {
     "payApp": "微信支付网关的应用名，如jiaomu-net",
     "out_bills_no": "商户系统内部的商家单号，只能由数字、大小写字母组成，在商户系统内部唯一"
  }
  ```

- 成功返回
  
  ```json
  {
     "code": 200,
     "msg": "OK",
     "result: {
        "out_bill_no" : "plfk2020042013",
        "transfer_bill_no" : "1330000071100999991182020050700019480001",
        "state" : "CANCELING",
        "update_time" : "2015-05-20T13:29:35.120+08:00"
     }
  }
  ```

#### 3. 商家转账查询

- 商家转账用户确认模式下，根据商户单号查询转账单的详细信息。

- 参考文档: https://pay.weixin.qq.com/doc/v3/merchant/4012716437

- 请求接口: POST /v3/query-transfer-bills

- Body:
  
  ```json
  {
     "payApp": "微信支付网关的应用名，如jiaomu-net",
     "out_bills_no": "商户系统内部的商家单号，只能由数字、大小写字母组成，在商户系统内部唯一"
  }
  ```

- 成功返回
  
  ```json
  {
     "code": 200,
     "msg": "OK",
     "result: {
        "mch_id" : "1900001109",
        "out_bill_no" : "plfk2020042013",
        "transfer_bill_no" : "1330000071100999991182020050700019480001",
        "appid" : "wxf636efh567hg4356",
        "state" : "SUCCESS",
        "transfer_amount" : 400000,
        "transfer_remark" : "新会员开通有礼",
        "fail_reason" : "PAYEE_ACCOUNT_ABNORMAL",
        "openid" : "o-MYE42l80oelYMDE34nYD456Xoy",
        "user_name" : "757b340b45ebef5467rter35gf464344v3542sdf4t6re4tb4f54ty45t4yyry45",
        "create_time" : "2015-05-20T13:29:35.120+08:00",
        "update_time" : "2015-05-20T13:29:35.120+08:00"
     }
  }
  ```

#### 4. 商家转账回调的验签和解密

- 微信支付系统通过商家转账回调通知接口通知商户系统单据处理到终态。

- 参考文档: https://pay.weixin.qq.com/doc/v3/merchant/4012712115

- 请求接口: POST /v3/verify-transfer-bills/:payApp

- 请求头：转发所有回调URL收到的请求头

- 请求体: 转发所有回调URL收到的请求体

- 失败返回
  
  ```json
  {
     "code": 500,
     "msg": "失败原因",
     "respToWxpay": { // 返回给微信支付服务的状态码和响应体
        "respCode": xxx,   // 返回给微信支付服务的状态码
        "respBody": {JSON} // 返回给微信支付服务的Body体
     }
  }
  ```

- 成功返回
  
  ```json
  {
     "code": 200,
     "msg": "OK",
     "respToWxpay": { // 返回给微信支付服务的状态码和响应体
        "respCode": xxx,   // 返回给微信支付服务的状态码
        "respBody": {JSON} // 返回给微信支付服务的Body体
     },
     "result": {
         "bill": {  // 解密后单据信息
            "out_bill_no": "plfk2020042013",
            "transfer_bill_no":"1330000071100999991182020050700019480001",
            "state": "SUCCESS",
            "mch_id": "1900001109",
            "transfer_amount": 2000,
            "openid": "o-MYE421800elYMDE34nYD456Xoy",
            "fail_reason:" "PAYEE_ACCOUNT_ABNORMAL",
            "create_time": "2015-05-20T13:29:35+08:00",
            "update_time": "2023-08-15T20:33:22+08:00"
         }
     }
  }
  ```
