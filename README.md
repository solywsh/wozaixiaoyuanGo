# wozaixiaoyuanGo

## 使用

在同一目录下创建`config.json`，并写入一下内容

| key           | function                         | type     |
| ------------- | -------------------------------- | -------- |
| jwsession     | 我在校园的密钥                   | string   |
| userAgent     | 设备标识                         | string   |
| pushPlusToken | PushPlus的token，用于接收名单    | []string |
| province      | 省份，日检日报提交哈希加密时使用 | string   |
| city          | 城市，日检日报提交哈希加密时使用 | string   |

demo：

```json
{
  "jwsession": "afbf1231231231231231319482131d",
  "userAgent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN",
  "province": "陕西省",
  "city": "西安市"
}

```



