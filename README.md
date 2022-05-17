# wozaixiaoyuanGo

## How to use

### 学生版（服务器后台运行）

创建或修改`config.yaml`文件

```yaml
user:
  # 姓名,id
  - name: 王世浩
    # 晨检
    morningCheck:
      enable: true # 是否开启,true/false,如果为空[checkTime]和[endTime]可以为空
      checkTime: '07:30:00' # 执行签到的时间
      endTime: '10:00:00' # 签到时间,签到的机制是,在大于[checkTime]并小于[endTime]之间执行,一天只会执行一次
    afternoonCheck:
      enable: true
      checkTime: '11:30:00'
      endTime: '15:00:00'
    eveningCheck:
      enable: true
      checkTime: '21:30:00'
      endTime: '23:59:00'
    # 我在校园密钥
    jwsession: ''

    # qq机器人,用做反馈
    # powered by: https://github.com/solywsh/qqBot-revue
    # 在使用前需要添加qq 3056159050为好友,发送`/help`根据提示获取密钥
    # 如果enable为false,[userId]和[token]可以为空
    # 也可以根据链接自己部署一个revue机器人,如果你有好的建议,欢迎到仓库提交issue
    qqBotRevue:
      enable: true
      userId: '' # 自己的qq
      token: '' # 获取到的token,注意qq和token为绑定关系
    # 设备标识,用以规避更换设备登录提示,这里建议使用自己手机的设备标识
    userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN'
    # 所在省份,在日检日报提交表单哈希时使用
    province: 陕西省
    # 所在省份的城市,在日检日报提交表单哈希时使用
    city: 西安市
```

如果你想同时给多人打卡：

```yaml
user:
  - name:
    morningCheck:
      enable: true 
      checkTime: 
      endTime: 
    afternoonCheck:
      enable: 
      checkTime: 
      endTime: 
    eveningCheck:
      enable: 
      checkTime: 
      endTime: 
    jwsession: 
    qqBotRevue:
      enable: 
      userId: 
      token: 
    userAgent: 
    province: 
    city: 
    ------------这里添加即可,注意对齐--------------
  - name:
    morningCheck:
      enable:
      checkTime: 
      endTime:
    afternoonCheck:
      enable:
      checkTime:
      endTime:
    eveningCheck:
      enable: true
      checkTime:
      endTime:
    jwsession:
    qqBotRevue:
      enable:
      userId:
      token:
    userAgent:
    province:
    city:
```

#### 注意

由于默认西安石油大学学生使用，在提交表单时并没有提供动态的更改表单内容，如果你为其他学校的学生，请执行更改`main.go`相关函数的内容。其中`CheckOperate()`用以执行日检日报，`doEveningCheck()`用以执行晚上定位签到。

### 导员版

创建或者修改`config.yaml`以下内容：

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
