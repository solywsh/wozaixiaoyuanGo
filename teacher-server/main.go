package main

import (
	"log"
	"time"
)

func operation() {
	dataNow := getDate()
	dateTmp := ""
	var yamlConfig Config
	var err error
	eventMap := make(map[string]map[string]int) // 记录今日任务执行flag
	log.Println("开始运行...")
	for {
		dataNow = getDate()
		// 第二天执行刷新
		if dataNow != dateTmp {
			dateTmp = getDate()
			yamlConfig, err = NewConf("./config.yaml")
			if err != nil {
				return // 读取错误退出
			}
			// 刷新flag,0为今日未执行
			for _, user := range yamlConfig.User {
				eventMap[user.Name] = map[string]int{"morning": 0}
				eventMap[user.Name] = map[string]int{"afternoon": 0}
				eventMap[user.Name] = map[string]int{"evening": 0}
				eventMap[user.Name] = map[string]int{"health": 0}
			}
		}
		timeNow := time.Now().Format("15:04:05")
		for _, user := range yamlConfig.User {
			if user.MorningCheck.Enable &&
				timeNow < user.MorningCheck.EndTime &&
				timeNow > user.MorningCheck.CheckTime &&
				eventMap[user.Name]["morning"] != 1 {
				eventMap[user.Name]["morning"] = 1 // flag 置为1
				// 晨检
				user.DailyCheck(1)
			}

			if user.AfternoonCheck.Enable &&
				timeNow < user.AfternoonCheck.EndTime &&
				timeNow > user.AfternoonCheck.CheckTime &&
				eventMap[user.Name]["afternoon"] != 1 {
				eventMap[user.Name]["afternoon"] = 1 // flag 置为1
				// 午检
				user.DailyCheck(2)
			}

			if user.EveningCheck.Enable &&
				timeNow < user.EveningCheck.EndTime &&
				timeNow > user.EveningCheck.CheckTime &&
				eventMap[user.Name]["evening"] != 1 {
				eventMap[user.Name]["evening"] = 1
				// 晚检
				user.EveningSignOperate()
			}

			if user.HealthCheck.Enable &&
				timeNow < user.HealthCheck.EndTime &&
				timeNow > user.HealthCheck.CheckTime &&
				eventMap[user.Name]["health"] != 1 {
				eventMap[user.Name]["health"] = 1
				// 晚检
				user.HealthCheckOperate()
			}
		}
	}
}

func main() {
	//operation()
	user := User{
		QqBotRevue: QqBotRevue{
			Enable: true,
			Module: "brief",
			Qq: []Qq{
				{
					UserId: "1228014966",
					Token:  "35406ce1-7dd2-41ca-9e46-d786deaca5f8",
				},
			},
		},
		Jwsession: "f46a2d76b23a449c841245fc9c245ec8",
		Name:      "王天航",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.25(0x1800192b) NetType/4G Language/zh_CN miniProgram/wx8a7eb7a1dbbba6cd",
	}

	user.DailyCheck(1)
}
