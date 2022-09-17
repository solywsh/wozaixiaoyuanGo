package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"log"
	"strconv"
	"time"
)

func getDate() string {
	return time.Now().Format("20060102")
}

func getSha256(src string) string {
	sha256Bytes := sha256.Sum256([]byte(src))
	sha256String := hex.EncodeToString(sha256Bytes[:])
	return sha256String
}

func (u User) qqBotRevueEvent(title, content string) {
	if !u.QqBotRevue.Enable {
		return
	}
	url := "http://revue.magicode123.cn:5000/send_private_msg"
	client := resty.New()
	_, err := client.R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
	}).SetBody(map[string]string{
		"token":   u.QqBotRevue.Token,
		"user_id": u.QqBotRevue.UserId,
		"message": "[" + title + "]" + "\n" + content,
	}).Post(url)
	if err != nil {
		return
	}
}

func (u User) CheckOperate(seq int) {
	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒
	content := u.Province + "_" + signTime + "_" + u.City
	signatureHeader := getSha256(content)
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetQueryParams(map[string]string{
		"answers":         "[\"0\"]",
		"seq":             strconv.Itoa(seq),
		"temperature":     "36.5",
		"userId":          "",
		"latitude":        "34.108216",
		"longitude":       "108.605084",
		"towncode":        "",
		"citycode":        "",
		"areacode":        "",
		"country":         "中国",
		"city":            "西安市",
		"district":        "鄠邑区",
		"province":        "陕西省",
		"township":        "甘亭街道",
		"street":          "东街",
		"myArea":          "610118",
		"timestampHeader": signTime,
		"signatureHeader": signatureHeader,
	}).Post("https://student.wozaixiaoyuan.com/heat/save.json")
	if err != nil {
		u.qqBotRevueEvent("日检日报提醒", "打卡失败，网络错误")
		log.Println(u.Name, "打卡失败，网络错误", "seq=", seq)
		return
	}
	var msg WzxyMessage
	err = json.Unmarshal(post.Body(), &msg)
	if err != nil {
		log.Println(u.Name, "打卡失败，解析错误", "seq=", seq)
	}
	fmt.Println(msg)
	postJson := gojsonq.New().JSONString(string(post.Body()))
	if int(postJson.Reset().Find("code").(float64)) == 0 {
		u.qqBotRevueEvent("日检日报提醒", "打卡成功")
		log.Println(u.Name, "打卡成功", "seq=", seq)
		// 正常
	} else {
		u.qqBotRevueEvent("日检日报提醒", "打卡失败,jwsession可能失效")
		log.Println(u.Name, "打卡失败,jwsession可能失效", "seq=", seq, post.String())
	}
}

type WzxyMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (u User) getSignMessage() (res int, signId, logId string) {
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"jwsession": u.Jwsession,
	}).SetQueryParams(map[string]string{
		"page": "1",
		"size": "5",
	}).Post("https://student.wozaixiaoyuan.com/sign/getSignMessage.json")
	if err != nil {
		return -1, "", ""
	}
	//fmt.Println(string(post.Body()))
	pJson := gojsonq.New().JSONString(string(post.Body()))
	if int(pJson.Reset().Find("code").(float64)) == 0 {
		signTimeStart := pJson.Reset().Find("data.[0].start")
		signTimeEnd := pJson.Reset().Find("data.[0].end")
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if timeNow > signTimeStart.(string) && timeNow < signTimeEnd.(string) {
			// 在签到区间
			signId = pJson.Reset().Find("data.[0].id").(string)
			logId = pJson.Reset().Find("data.[0].logId").(string)
			return 0, signId, logId
		} else {
			// 在签到区间
			signId = pJson.Reset().Find("data.[0].id").(string)
			logId = pJson.Reset().Find("data.[0].logId").(string)
			//fmt.Println(signId, logId)
			// 不在签到区间
			return -2, "", ""
		}
	}
	return -1, "", ""
}

func (u User) doEveningCheck(signId, logId string) {
	url := "https://student.wozaixiaoyuan.com/sign/doSign.json"
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetBody(map[string]string{
		"signId":    signId,
		"towncode":  "",
		"citycode":  "",
		"areacode":  "",
		"city":      "西安市",
		"id":        logId,
		"latitude":  "34.10154079861111",
		"longitude": "108.65831163194444",
		"country":   "中国",
		"district":  "鄠邑区",
		"township":  "五竹街道",
		"province":  "陕西省",
	}).Post(url)
	if err != nil {
		return
	}
	pJson := gojsonq.New().JSONString(string(post.Body()))
	if int(pJson.Reset().Find("code").(float64)) == 0 {
		u.qqBotRevueEvent("晚检签到提醒", "晚检签到成功")
		log.Println(u.Name, "晚检签到成功")
	} else {
		u.qqBotRevueEvent("晚检签到提醒", "晚检签到失败,返回信息为:"+string(post.Body()))
		log.Println(u.Name, "晚检签到失败,返回信息为:", string(post.Body()))
	}
}

func (u User) EveningCheckOperate() {
	res, signId, logId := u.getSignMessage()
	switch res {
	case 0:
		// 正常,执行签到
		// fmt.Println(signId, logId)
		u.doEveningCheck(signId, logId)
	case -1:
		u.qqBotRevueEvent("晚检签到提醒", "获取晚检信息失败,网络错误")
		log.Println(u.Name, "获取晚检信息失败,网络错误")
	case -2:
		u.qqBotRevueEvent("晚检签到提醒", "晚检签到失败,不在签到时间范围内")
		log.Println(u.Name, "晚检签到失败,不在签到时间范围内")
	}
}

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
				go user.CheckOperate(1)
			}

			if user.AfternoonCheck.Enable &&
				timeNow < user.AfternoonCheck.EndTime &&
				timeNow > user.AfternoonCheck.CheckTime &&
				eventMap[user.Name]["afternoon"] != 1 {
				eventMap[user.Name]["afternoon"] = 1 // flag 置为1
				// 午检
				go user.CheckOperate(2)
			}

			if user.EveningCheck.Enable &&
				timeNow < user.EveningCheck.EndTime &&
				timeNow > user.EveningCheck.CheckTime &&
				eventMap[user.Name]["evening"] != 1 {
				eventMap[user.Name]["evening"] = 1
				// 晚检
				go user.EveningCheckOperate()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	// operation()
	user := User{
		Jwsession: "1bfa5832eb3e406e9a3d73cc882eda13",
		QqBotRevue: QqBotRevue{
			Enable: true,
			Token:  "35406ce1-7dd2-41ca-9e46-d786deaca5f8",
			UserId: "1228014966",
		},
		Name:           "王世浩",
		MorningCheck:   MorningCheck{},
		AfternoonCheck: AfternoonCheck{},
		EveningCheck:   EveningCheck{},
		Province:       "陕西省",
		City:           "西安市",
		UserAgent:      "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN",
	}
	user.CheckOperate(1)
}
