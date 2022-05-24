package main

import (
	"time"
)

func operation() {
	dataNow := getDate()
	dateTmp := ""
	var yamlConfig Config
	var err error
	eventMap := make(map[string]map[string]int) // 记录今日任务执行flag

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
				go user.DailyCheck(1)
			}

			if user.AfternoonCheck.Enable &&
				timeNow < user.AfternoonCheck.EndTime &&
				timeNow > user.AfternoonCheck.CheckTime &&
				eventMap[user.Name]["afternoon"] != 1 {
				eventMap[user.Name]["afternoon"] = 1 // flag 置为1
				// 午检
				go user.DailyCheck(2)
			}

			if user.EveningCheck.Enable &&
				timeNow < user.EveningCheck.EndTime &&
				timeNow > user.EveningCheck.CheckTime &&
				eventMap[user.Name]["evening"] != 1 {
				eventMap[user.Name]["evening"] = 1
				// 晚检
				go user.EveningSignOperate()
			}
		}
	}
}

func main() {
	//operation()
	yamlConfig, _ := NewConf("./config.yaml")
	yamlConfig.User[0].qqBotRevueEvent("标题", "人数为:", []string{"王世浩", "康宇"})
	yamlConfig.User[0].QqBotRevue.Module = "detailed"
	yamlConfig.User[0].qqBotRevueEvent("标题", "人数为:", []string{"王世浩", "康宇"})
}
