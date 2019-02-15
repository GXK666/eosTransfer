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
