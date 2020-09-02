你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# eosTransfer


## 一、钱包转出api：

> 服务程序管理eos私钥，用于签名交易。 客户端程序调用api即可发出转账

```
# post json 格式参数

post   /v1/transfer_out

{
  "contract":"eosio.token",  # 转eos代币，则此参数为eosio.token
  "to:"accountB",
  "amount":"1.0000 EOS",     # 代币标准格式
  "memo": "memo"
  "request_id":"uuid"
}


# api返回结果
{
    "txid":""  # 64位交易id，后续查询交易是否成功需要用到
}

```


## 二、查询交易


```
post /v1/get_transfer"

{
    "txid": ""  # 64位交易id
}


# 返回结果:irreversible (打包成功且不可逆，大约在交易发出3分钟后查询进入此状态)， executed（交易打包成功但可逆）
{
   "status":"irreversible"  # 'irreversible' 'executed','soft_fail','hard_fail','delayed','expired','unknown'
}
```
