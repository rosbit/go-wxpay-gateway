# 通用微信支付网关服务

1. `go-wxpay-gateway`是封装了与微信支付的相关接口，以更简便易用的形式(封装了各种编码、加解密、证书、
    随机数、沙箱的使用)让应用使用微信支付功能
2. `go-wxpay-gateway`运行后，通过HTTP给需要微信支付的应用提供接口，请求与响应的数据都是JSON，任何支持
    HTTP的语言都可以实现业务代码。

## 架构图

 ![架构图](wxpay-gateway.png)

- 任何商户平台的配置信息都统一由`go-wxpay-gateway`集中管理，应用只要通过`应用名`、`appId`就可以使用这些参数
- `go-wxpay-gateway`提供接口解析支付/退款通知程序接收的参数，并把解析后的结果、返回给微信支付的结果返回接口调用程序
- 下单、支付、支付结果通知流程
  1. `微信支付应用`准备好支付信息(订单号、支付金额)，根据应用类型选取相应的支付方式(JSAPI/H5/APP/NATIVE)，
      同时把`应用名`、`appId`、`支付结果通知URL`提供给`go-wxpay-gateway`
  2. `go-wxpay-gateway`取出`应用名`对应的配置信息，向`微信商户平台`发出统一支付接口调用
  3. `微信商户平台`返回相应的支付参数
  4. `go-wxpay-gateway`抽取支付相关的参数，**生成**或**包装**`激活微信支付`的参数，返回给`微信支付应用`
  5. `微信支付应用`把相关参数返回给前端 或者 提供二维码 激活`微信app`
  6. 用户通过`微信app`完成支付
  7. `微信商户平台`异步把支付结果通知给`支付结果通知URL`，该URL由应用自己完成接收
  8. `支付结果通知URL`处理程序把收到的信息和`应用名`转发给`go-wxpay-gateway`提供的参数解析接口
  9. `go-wxpay-gateway`根据其中的`应用名`找到相应配置参数，对支付结果通知参数进行验证解析
  10. `go-wxpay-gateway`把`结果通知`的解析结果、返回给微信支付的结果返回给接口调用程序

## 下载、编译方法

1. 前提：已经安装go 1.11.x及以上、git、make

2. 进入任一文件夹，执行命令
   
   ```bash
   $ git clone https://github.com/rosbit/go-wxpay-gateway
   $ cd go-wxpay-gateway
   $ make
   ```

3. 编译成功，会得到3个可执行程序
   
   - `go-wxpay-gateway`: 微信支付网关程序，可以执行`./go-wxpay-gateway -v`显示程序信息。
   - `go-wxpay-getsandbox`: 获取沙箱测试的apiKey工具，可以执行`./go-wxpay-getsandbox -v`显示程序信息。

4. Linux的二进制版本可以直接进入[releases](https://github.com/rosbit/go-wxpay-gateway/releases)下载

## 运行方法

### go-wxpay-gateway运行方法

1. 环境变量
   - CONF_FILE: 指明配置的路径，格式见下文
2. 配置文件格式
   - 是一个JSON
   - 可以通过`wxpay-gateway.conf.sample.json`进行修改:
   - 其中的`merchants`定义了商户的信息，实际账号和沙箱账号请分别配置
   - 其中的`apps`定义了应用的信息，项`name`就是`应用名`
     - `应用名`只要可以区分`微信支付应用`就可以了
     - 如果使用沙箱测试，另外配置一个带`-dev`后缀的`应用名`，如`app-dev`，并指向沙箱merchant信息，
       `微信支付应用`使用`app-dev`时不产生实际的费用；等测试成功后，再切换到实际的`应用名`，如`app`
   - 配置项`endpoints`中的路由确保内网可以使用就可以了
3. 运行
   - `$ CONF_FILE=./wxpay-gateway-conf.json ./go-wxpay-gateway`

### go-wxpay-getsandbox运行方法

- 运行`./go-wxpay-getsandbox <appId> <mchId> <mchApiKey>`
  - 其中&lt;appId&gt;是指公众号、小程序、网页应用等id，请根据应用类型选用正确的appId
  - &lt;mchId&gt;是指商户号，申请商户号的时可以拿到
  - &lt;mchApiKey&gt;是商户号的apiKey，可以进入商户号管理界面设定
- 运行成功后把沙箱mchApiKey输出到屏幕

## go-wxpay-gateway服务的接口参数

1. 创建订单
   
   - 对应配置项: `create-pay`
   
   - 访问方法: POST
   
   - URI: 如配置值为`/wxpay/create-pay`，则通过`/wxpay/create-pay/<type>`访问，&lt;type&gt;的可取值为`JSAPI`、`H5`、`APP`、`NATIVE`
   
   - **[注意]** 由于该接口会使用商户相关的配置信息，该接口不要暴露在外网，只要支付应用可以访问就可以了
   
   - 以下为各种支付方式的参数、响应结果说明:
     
     1. JSAPI支付
        
        - URI: `/wxpay/create-pay/JSAPI`
        
        - 请求参数
          
          ```javascript
          {
             "appId": "公众号或小程序的appId",
             "payApp": "go-wxpay-gateway配置文件中的应用名",
             "goods": "商品名",
             "udd": "用户自定义数据",
             "orderId": "支付应用中唯一的订单号",
             "fee": 以分为单位的支付额度,
             "ip": "创建订单的IP地址",
             "openId": "支付用户在公众号或小程序中的openId",
             "notifyUrl": "支付结果通知URL，必须外网可以访问"
          }
          ```
        
        - 响应结果
          
          ```json
          {
              "code": 200,
              "msg": "OK",
              "result":{
                 "jsapi_params":{
                    "appId": "wxc914c55fba939e3c",
                    "nonceStr": "Q6cU6oyVVvSsIWusb20dett1X59h4ezA",
                    "package": "prepay_id=wx20190530150713356368",
                    "paySign": "4F676F3FF2F4719289EA0685928CF83B",
                    "signType": "MD5",
                    "timeStamp": "1559200033"
                 },
                 "prepay_id": "wx20190530150713356368"
               }
          }
          ```
     
     2. H5支付
        
        - URI: `/wxpay/create-pay/H5`
        
        - 请求参数
          
          ```json
          {
              "appId": "公众号或小程序的appId",
              "payApp": "go-wxpay-gateway配置文件中的应用名",
              "goods": "商品名",
              "udd": "用户自定义数据",
              "orderId": "支付应用中唯一的订单号",
              "fee": 以分为单位的支付额度,
              "ip": "创建订单的IP地址",
              "redirectUrl": "在支付成功后需要跳转的URL",
              "notifyUrl": "支付结果通知URL，必须外网可以访问"
              "sceneInfo": {
                 "h5_info": {
                   "type": "类型",
                   "wap_name": "wap应用名",
                   "wap_url": "wap-site-url"
                 }
                 // ---- OR ----
                 "h5_info": {
                   "type": "类型",
                   "app_name":"ios应用名",
                   "bundle_id": "ios-bundle-id"
                 }
                 // ---- OR ----
                 "h5_info": {
                   "type": "类型",
                   "app_name":"android应用名",
                   "package_name": "android-package-name"
                 }
               }
          }
          ```
        
        - 响应结果
          
          ```json
          {
               "code": 200,
               "msg": "OK",
               "result":{
                   "pay_url": "weixin://wxpay/s/An4baqw?redirect_url=http%3A%2F%2Flocalhost%3A11083",
                   "prepay_id": "wx20190530150945625582"
                }
          }
          ```
     
     3. APP支付
        
        - URI: `/wxpay/create-pay/APP`
        
        - 请求参数
          
          ```json
          {
              "appId": "公众号或小程序的appId",
              "payApp": "go-wxpay-gateway配置文件中的应用名",
              "goods": "商品名",
              "udd": "用户自定义数据",
              "orderId": "支付应用中唯一的订单号",
              "fee": 以分为单位的支付额度,
              "ip": "创建订单的IP地址",
              "notifyUrl": "支付结果通知URL，必须外网可以访问"
          }
          ```
        
        - 响应结果
          
          ```json
          {
              "code": 200,
              "msg": "OK",
              "result":{
                 "prepay_id": "wx20190530151116852853",
                 "req_params":{
                    "appid": "wxc914c55fba939e3c",
                    "noncestr": "iAQIoaqU854JHNnDhG00C0QSbTdIfthD",
                    "package": "Sign=WXPay",
                    "partnerid": "Iamyourpartner",
                    "prepayid": "wx20190530151116852853",
                    "sign": "1ADDF00C546912C446D9B9EA847B0D0F",
                    "timestamp": "1559200276"
                 }
               }
          }
          ```
     
     4. NATIVE支付
        
        - URI: `/wxpay/create-pay/NATIVE`
        
        - 请求参数
          
          ```json
          {
              "appId": "公众号或小程序的appId",
              "payApp": "go-wxpay-gateway配置文件中的应用名",
              "goods": "商品名",
              "udd": "用户自定义数据",
              "orderId": "支付应用中唯一的订单号",
              "fee": 以分为单位的支付额度,
              "ip": "创建订单的IP地址",
              "productId": "商品ID",
              "notifyUrl": "支付结果通知URL，必须外网可以访问"
           }
          ```
        
        - 响应结果
          
          ```json
          {
              "code": 200,
              "msg": "OK",
              "result":{
                  "code_url": "weixin://wxpay/s/An4baqw",
                  "prepay_id": "wx20190530151411214102"
               }
          }
          ```

2. 查询订单
   
   - 对应配置项: `query-order`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/wxpay/query-order`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 请求参数
     
     ```json
     {
         "appId": "应用的appId",
         "payApp": "go-wxpay-gateway配置文件中的应用名",
         "orderId": "支付应用中的唯一订单id号"
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "app_id": "wxc914c55fba939e3c",
         "mch_id": "1530730681",
         "device_info": "sandbox",
         "result_code": "SUCCESS",
         "err_code": "SUCCESS",
         "err_code_des": "SUCCESS",
         "open_id": "wxd930ea5d5a258f4f",
         "is_subscribe": true,
         "trade_type": "JSAPI",
         "bank_type": "CMC",
         "total_fee": 101,
         "settlement_total_fee": 101,
         "fee_type": "CNY",
         "cash_fee": 101,
         "cash_fee_type": "CNY",
         "coupon_fee": 0,
         "coupon_count": 0,
         "coupons": null,
         "transaction_id": "4662714807620190530105704209783",
         "order_id": "o002",
         "attach": "sandbox_attach",
         "time_end": "20190530105704"
     }
     ```

3. 关闭订单
   
   - 对应配置项: `close-order`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/wxpay/close-order`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 请求参数
     
     ```json
     {
         "appId": "应用的appId",
         "payApp": "go-wxpay-gateway配置文件中的应用名",
         "orderId": "支付应用中的唯一订单id号"
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "code": 200,
         "msg": "OK"
     }
     ```

4. 创建退款
   
   - 对应配置项: `create-refund`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/wxpay/create-refund`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 请求参数
     
     ```json
     {
         "appId": "应用的appId",
         "payApp": "go-wxpay-gateway配置文件中的应用名",
         "orderId": "支付应用中的唯一订单id号",
         "refundId": "支付应用中的唯一退款id号",
         "totalFee": 支付的总费用，单位分,
         "refundFee": 退款的费用，单位分,
         "refundReason": "退款理由",
         "notifyUrl": "支付结果通知URL，必须外网可以访问"
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "code": 200,
         "msg": "OK",
         "result":{
           "app_id": "wxc914c55fba939e3c",
           "mch_id": "1530730681",
           "result_code": "",
           "err_code": "SUCCESS",
           "err_code_des": "",
           "transaction_id": "4630338607020190530111333587924",
           "order_id": "o002",
           "wx_refund_id": "4630338607020190530111333587",
           "refund_id": "r002",
           "total_fee": 101,
           "settlement_total_fee": 0,
           "refund_fee": 101,
           "settlement_refund_fee": 0,
           "fee_type": "CNY",
           "cash_fee": 101,
           "cash_fee_type": "CNY",
           "cash_refund_fee": 101,
           "coupon_refund_fee": 0,
           "coupon_refund_count": 0,
           "refund_coupons": null
         }
     }
     ```

5. 校验通知结果
   
   - 校验支付支付结果
     
     - URI: 对应配置文件里`endpoints`列表中的`verify-notify-pay` + '?app=<应用名称>'
     
     - 方法: POST
     
     - 请求参数: 应用支付结果URL收到的POST Body内容
     
     - 响应结果:
       
       ```json
       {
          "code": 200/406, // 正确返回200，失败返回406
          "msg": "OK/解析失败原因",
          "params": null, // code为406时返回
          "params": {     // code为200时返回，是从Post Body中解析出来的
             "result_code": "SUCCESS",
             "err_code": "SUCCESS",
             "err_code_des": "SUCCESS",
             "mch_id": "1530730681",
             "device_info": "APP",
             "trade_type": "APP",
             "attach": "any",
             "bank_type": "CMC",
             "fee_type": "CNY",
             "total_fee": 101,
             "settlement_total_fee": 101,
             "cash_fee_type": "CNY",
             "cash_fee": 101,
             "coupon_count": 0,
             "coupons": null,
             "is_subscribe": true,
             "open_id": "sandboxopenid",
             "order_id": "o003",
             "time_end": "20190530151118",
             "transaction_id": "4391741731420190530151118238943"
          },
          "msgForWxpay": { // 返回给微信支付服务的内容，是一个XML
              "<xml>....</xml>"
          }
       }
       ```
   
   - 校验退款结果
     
     - URI: 对应配置文件里`endpoints`列表中的`verify-notify-refund` + '?app=<应用名>'
     - 方法: POST
     - 请求参数: 应用退款结果URL收到的POST Body内容
     - 响应结果:
       - 参考校验支付支付结果

6. V3版商家转账到零钱
   
   - 对应配置项: `v3-transfer`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/v3/transfer`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 本接口只是“发出转账申请”，成功申请后需要在后台审核确认才能转出零钱
   
   - 请求参数
     
     ```json
     {
         "payApp": "go-wxpay-gateway配置文件中的应用名",
         "appId": "应用的appId",
         "batchNo": "本次转账的批次，应用内唯一",
         "batchName": "批次名称",
         "batchRemark": "批次描述",
         "details": [
            {
              "tradeNo": "批次内唯一交易id",
              "amount": 30, // 金额，单位"分",
              "desc": "描述",
              "openId": "收款人在appId内的openId",
              "userName": "2000以上，必须给收款人实名"
            },
            {...}
         ]
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "code": 200,
         "msg": "OK",
         "result":{
            "wxBatchId":"1030000040101068797782022081801027450232" // 微信批次id
         }
     }
     ```

7. V3版商家批次查询
   
   - 对应配置项: `v3-query-transfer`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/v3/query-transfer`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 请求参数
     
     ```json
     {
         "payApp": "go-wxpay-gateway配置文件中的应用名",
         "status": "查询的状态值",  //  "ALL" | "SUCCESS" | "FAIL", 不给用ALL
         "needDetail": true, // 是否展示详情
         "wxBatchId": "转账申请成功返回的微信批次id"
            "------ or -----": "或者",
         "batchNo": "转账的批次，应用内唯一"
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "code": 200,
         "msg": "OK",
         "result":{
            "status":"FINISHED",
            "total": 1, // 数目
            "details": [
               {
                  "detail_id":"1040000040101068797782022081801020503999", // 微信详情id
                  "out_detail_no":"tt202208170001", // 对应自己的tradeNo
                  "detail_status":"SUCCESS"         // 详情状态
               }
            ]
         }
     }
     ```

8. V3版商家转账详情查询
   
   - 对应配置项: `v3-query-transfer-detail`
   
   - 访问方法: POST
   
   - URI: 直接根据配置值访问，如`/v3/query-transfer-detail`
   
   - **[注意]** 该接口配置成一个内网可以访问的API
   
   - 请求参数
     
     ```json
     {
         "payApp": "go-wxpay-gateway配置文件中的应用名",
     
         "wxBatchId": "转账申请成功返回的微信批次id",
         "wxDetailId": "批次查询时返回的微信详情id",
            "------ or -----": "或者",
         "batchNo": "转账的批次，应用内唯一",
         "tradeNo": "转账时详情的唯一id"
     }
     ```
   
   - 响应结果
     
     ```json
     {
         "code": 200,
         "msg": "OK",
         "result":{
             "status":"SUCCESS",  // 详情状态
             "amount": 30, // 金额,单位分
             "remark": "转账时的备注",
             "reason": "失败原因。如果成功，返回空"
         }
     }
     ```
