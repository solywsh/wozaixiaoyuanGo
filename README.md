# wozaixiaoyuanGo

## How to use

修改或添加`config.yaml`以下内容：

| key       | function                         |
| --------- | -------------------------------- |
| jwsession | 我在校园的密钥，抓包获得         |
| userAgent | 设备标识，抓包获得               |
| province  | 省份，日检日报提交哈希加密时使用 |
| city      | 城市，日检日报提交哈希加密时使用 |

比如`config.yaml`：

```yaml
# 我在校园密钥,抓包获得
jwsession: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# 用户设备的标识
userAgent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN

# 所在省份和城市
province: 陕西省
city: 西安市
```

然后执行程序即可：

![demo](demo.gif)
