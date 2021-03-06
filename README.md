# wozaixiaoyuanGo

## How to use

### 学生版-服务器运行

> 此版本建议配置完成之后运行在服务器端

![QQ图片20220607163750](http://cdnimg.violetwsh.com/img/QQ%E5%9B%BE%E7%89%8720220607163750.png)

在`wozaixiaoyuanGo/student`目录下，创建或修改`config.yaml`文件（注意对齐否则可能无法识别）

```yaml
user:
  # 姓名,id
  - name: '张三'
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

> 注意

由于默认西安石油大学学生使用，在提交表单时并没有提供动态的更改表单内容，如果你为其他学校的学生，请执行更改`main.go`相关函数的内容。其中`CheckOperate()`用以执行日检日报，`doEveningCheck()`用以执行晚上定位签到。

### 老师版-客户端运行

> 此版本可以直接在电脑上运行，如果你在Windows系统上运行，下载[Releases · solywsh/wozaixiaoyuanGo (github.com)](https://github.com/solywsh/wozaixiaoyuanGo/releases)的压缩包完成配置文件即可使用。

在`wozaixiaoyuanGo/teacher`目录下，创建或者修改`config.yaml`以下内容：

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

然后执行程序即可（由于后期改为goroutine的方式发送数据，速度可能会快一些），示例：

![demo](demo.gif)

### 老师版-服务器运行

> 此版本建议配置完成之后运行在服务器端，并且支持多人。

![img](http://cdnimg.violetwsh.com/img/%5D@E1%5D4%7D$~LG%7D%7DC@TZ_~K_~B.png)

在`wozaixiaoyuanGo/teacher-server`目录下，创建或者修改`config.yaml`以下内容：

```yaml
user:
  # 姓名,id
  - name: '张三'
    # 晨检
    morningCheck:
      enable: true # 是否开启,true/false,如果为空[checkTime]和[endTime]可以为空
      checkTime: '09:50:00' # 执行签到的时间
      endTime: '10:00:00' # 签到时间,签到的机制是,在大于[checkTime]并小于[endTime]之间执行,一天只会执行一次
    # 午检
    afternoonCheck:
      enable: true
      checkTime: '14:50:00'
      endTime: '15:00:00'
    # 晚检签到(可直接用于签到)
    eveningCheck:
      enable: true
      checkTime: '23:50:00'
      endTime: '23:59:00'
    # 健康打卡(西石大研究生使用)
    healthCheck:
      enable: false
      checkTime: '12:00:00'
      endTime: '23:59:00'
    # 我在校园密钥
    jwsession: ''
    # qq机器人,用做反馈和接收代签情况
    # powered by: https://github.com/solywsh/qqBot-revue
    # 在使用前需要添加qq 3056159050为好友,发送`/help`根据提示获取密钥
    # 如果enable为false,[userId]和[token]可以为空
    # 也可以根据链接自己部署一个revue机器人,如果你有好的建议,欢迎到仓库提交issue
    qqBotRevue:
      enable: true # 是否开启true/false
      module: brief # 消息发送方式brief/detailed,brief只发送人数,detailed发送名单
      qq:
        - userId: '' # 接收消息的qq1
          token: '' # 获取到的token,注意qq和token为绑定关系
        - userId: '' # 接收消息的qq2
          token: '' # 获取到的token,注意qq和token为绑定关系
    # 设备标识,用以规避更换设备登录提示,这里建议使用自己手机的设备标识
    userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN'
    # 所在省份,在日检日报提交表单哈希时使用
    province: '陕西省'
    # 所在省份的城市,在日检日报提交表单哈希时使用
    city: '西安市'
  - name: '李四'
    # 晨检
    morningCheck:
      enable: true
      checkTime: '09:30:00'
      endTime: '10:00:00'
    afternoonCheck:
      enable: true
      checkTime: '14:30:00'
      endTime: '15:00:00'
    healthCheck:
      enable: false
      checkTime: '12:00:00'
      endTime: '23:59:00'
    eveningCheck:
      enable: false
      checkTime: '21:30:00'
      endTime: '23:59:00'
    jwsession: ''
    qqBotRevue:
      enable: true
      module: brief
      qq:
        - userId: ''
          token: ''
        - userId: ''
          token: ''
    userAgent: ''
    province: ''
    city: ''
```

