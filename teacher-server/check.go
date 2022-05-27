package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"log"
	"strconv"
	"sync"
	"time"
)

// arg[0]:title
// arg[1]:content
// arg[2]:nameList []string
//
//example:
//  yamlConfig, _ := NewConf("./config.yaml")
//	yamlConfig.User[0].qqBotRevueEvent("标题", "内容")
//	yamlConfig.User[0].qqBotRevueEvent("标题", "人数为:", []string{"张三", "李四"})
//	yamlConfig.User[0].QqBotRevue.Module = "detailed"
//	yamlConfig.User[0].qqBotRevueEvent("标题", "未签到名单:", []string{"张三", "李四"})
//result:
//	[标题]
//	内容
//	[标题]
//	人数为:2
//	[标题]
//	未签到名单:张三 李四
func (u User) qqBotRevueEvent(arg ...interface{}) {
	if !u.QqBotRevue.Enable {
		return
	}
	var content string
	if len(arg) >= 3 {
		if u.QqBotRevue.Module == "brief" {
			content = arg[1].(string) + strconv.Itoa(len(arg[2].([]string)))
		} else {
			content = arg[1].(string)
			for _, name := range arg[2].([]string) {
				content += name + " "
			}
		}
	} else {
		content = arg[1].(string)
	}

	url := "http://revue.magicode123.cn:5000/send_private_msg"
	// 发送消息列表
	for _, qq := range u.QqBotRevue.Qq {
		go func(qq Qq) {
			client := resty.New()
			_, err := client.R().SetHeaders(map[string]string{
				"Content-Type": "application/json",
			}).SetBody(map[string]string{
				"token":   qq.Token,
				"user_id": qq.UserId,
				"message": "[" + arg[0].(string) + "]" + "\n" + content,
			}).Post(url)
			if err != nil {
				log.Println("[qqBotRevueEvent]", u.Name, "发送消息失败:", err.Error())
			}
		}(qq)
	}
}

var wg sync.WaitGroup

func getDate() string {
	return time.Now().Format("20060102")
}

func getSha256(src string) string {
	sha256Bytes := sha256.Sum256([]byte(src))
	sha256String := hex.EncodeToString(sha256Bytes[:])
	return sha256String
}

func (u User) checkForStudent(stuName, stuId, seq string) {
	defer wg.Done()

	province := u.Province
	city := u.City
	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒

	content := province + "_" + signTime + "_" + city
	signatureHeader := getSha256(content)
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetQueryParams(map[string]string{
		"answers":         "[\"0\"]",
		"seq":             seq,
		"temperature":     "36.5",
		"userId":          stuId,
		"latitude":        "",
		"longitude":       "",
		"country":         "",
		"city":            "",
		"district":        "",
		"province":        "",
		"township":        "",
		"street":          "",
		"areacode":        "",
		"timestampHeader": signTime,
		"signatureHeader": signatureHeader,
	}).Post("https://teacher.wozaixiaoyuan.com/heat/save.json")
	if err != nil {
		log.Println("[checkForStudent]", u.Name, stuName, "执行代打卡发生错误:"+err.Error())
		return
	}
	postJson := gojsonq.New().JSONString(string(post.Body()))
	if int(postJson.Reset().Find("code").(float64)) != 0 {
		log.Println("[checkForStudent]", u.Name, stuName, "执行代打卡请求错误:"+string(post.Body()))
		return
	}
	// log.Println("[doSignEvening]", u.Name, stuName, "签到成功")

}

func (u User) DailyCheck(seq int) {
	log.Println("[dailyCheck]seq =", seq, u.Name, "开始代签")
	client := resty.New()
	page := 1
	var unsignedStuId []string
	var unsignedName []string
	for {
		post, err := client.R().SetQueryParams(map[string]string{
			"seq":  strconv.Itoa(seq),
			"date": getDate(),
			"type": "0",
			"page": strconv.Itoa(page),
			"size": "20",
		}).SetHeaders(map[string]string{
			"JWSESSION":  u.Jwsession,
			"User-Agent": u.UserAgent,
		}).Post("https://teacher.wozaixiaoyuan.com/heat/getHeatUsers.json")
		if err != nil {
			log.Println("[dailyCheck]", u.Name, "未打卡名单请求错误,错误信息为:"+err.Error())
			return
		}
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		if postInfo.Find("code") != -10 {
			if len(postInfo.Reset().Find("data").([]interface{})) == 0 {
				if page == 1 {
					log.Println("[dailyCheck]seq=", seq, u.Name, "没有打卡信息或者打卡没有开始!")
					u.qqBotRevueEvent("日检日报代签提醒", "没有打卡信息或者打卡没有开始!")
					return
				}
				break
			}
			unsignedData := postInfo.Reset().From("data").Select("userId", "name").Get()
			for _, data := range unsignedData.([]interface{}) {
				unsignedStuId = append(unsignedStuId, data.(map[string]interface{})["userId"].(string))
				unsignedName = append(unsignedName, data.(map[string]interface{})["name"].(string))
			}
			page++
			//time.Sleep(1 * time.Second)
		} else {
			log.Println("[dailyCheck]seq=", seq, u.Name, "jwsession失效,请更换!")
			u.qqBotRevueEvent("日检日报代签提醒", "jwsession失效,请更换!")
			break
		}
	}
	time.Sleep(1 * time.Second)
	wg.Add(len(unsignedStuId))
	for i := 0; i < len(unsignedStuId); i++ {
		go u.checkForStudent(unsignedName[i], unsignedStuId[i], strconv.Itoa(seq))
		//time.Sleep(1 * time.Second)
	}
	wg.Wait()
	log.Println("[dailyCheck]seq =", seq, u.Name, "打卡完成!")
	if u.QqBotRevue.Module == "brief" {
		u.qqBotRevueEvent("日检日报代签提醒", "("+u.Name+")"+"代签人数为:", unsignedName)
	} else {
		u.qqBotRevueEvent("日检日报代签提醒", "("+u.Name+")"+"代签名单为:", unsignedName)
	}

}

func (u User) getEveningSignId() (signId string) {
	client := resty.New()

	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION": u.Jwsession,
	}).SetQueryParams(map[string]string{
		"state":   "0",
		"keyword": "",
		"page":    "1",
	}).Post("https://teacher.wozaixiaoyuan.com/signManage/getList.json")
	if err != nil {
		log.Println("[getEveningSignId]", u.Name, "请求晚检签到id发生错误,错误信息为:"+err.Error())
		return "0" //发生错误
	}
	//fmt.Println(string(post.Body()))
	postInfo := gojsonq.New().JSONString(string(post.Body()))
	if int(postInfo.Find("code").(float64)) == 0 {
		signEndTime := postInfo.Reset().Find("data.[0].end").(string)
		signBeginTime := postInfo.Reset().Find("data.[0].start").(string)
		signId := postInfo.Reset().Find("data.[0].id").(string)
		signNowTime := time.Now().Format("2006-01-02 15:04")
		if signNowTime > signBeginTime && signEndTime > signNowTime {
			//fmt.Println(signId)
			return signId // 正常返回签到Id
		} else {
			return "1" // 未到(或已过)签到时间
		}
	}
	return "0"
}

func (u User) getUnsignedList(signId string) (unsignedList []map[string]interface{}) {

	url := "https://teacher.wozaixiaoyuan.com/signManage/getSignLogs.json"
	client := resty.New()
	page := 1
	for {
		post, err := client.R().SetHeaders(map[string]string{
			"JWSESSION": u.Jwsession,
		}).SetQueryParams(map[string]string{
			"id":       signId,
			"page":     strconv.Itoa(page),
			"type":     "0",
			"size":     "20",
			"targetId": "",
		}).Post(url)
		if err != nil {
			log.Println("[getEveningSignId]", u.Name, "请求晚检未签到名单发生错误,错误信息为:"+err.Error())
			return unsignedList
		}
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		//fmt.Println(reflect.TypeOf(postInfo.Find("data.[0]")))
		if int(postInfo.Reset().Find("code").(float64)) == 0 && len(postInfo.Reset().Find("data").([]interface{})) != 0 {
			selectInfo := postInfo.Reset().From("data").Select("name", "id").Get() // 只提取出id和name
			for _, info := range selectInfo.([]interface{}) {
				unsignedList = append(unsignedList, info.(map[string]interface{}))
				//fmt.Println(info) // map[id:437356871436734465 name:马鹏]
			}
		} else {
			break
		}
		page++
		//time.Sleep(1 * time.Second)
	}

	return unsignedList
}

func (u User) doSignEvening(unsignedList []map[string]interface{}) {
	log.Println("[doSignEvening]", u.Name, "开始代签")
	url := "https://teacher.wozaixiaoyuan.com/signManage/adminSign.json"
	var unsignedName []string
	wg.Add(len(unsignedList))
	for _, unsignedInfo := range unsignedList {
		unsignedName = append(unsignedName, unsignedInfo["name"].(string))
		go func(info map[string]interface{}) {
			defer wg.Done()

			client := resty.New()
			post, err := client.R().SetHeaders(map[string]string{
				"JWSESSION": u.Jwsession,
			}).SetQueryParams(map[string]string{
				"id":   info["id"].(string),
				"type": "1",
			}).Post(url)
			if err != nil {
				log.Println("[doSignEvening]", u.Name, info["name"].(string), "执行晚检代签发生错误,错误信息为:"+err.Error())
				return
			}
			rJson := gojsonq.New().JSONString(string(post.Body()))
			if int(rJson.Reset().Find("code").(float64)) != 0 {
				log.Println("[doSignEvening]", u.Name, info["name"].(string), "代签失败,返回信息为:"+string(post.Body()))
			} else {
				//log.Println("[doSignEvening]", u.Name, info["name"].(string), "签到成功")
			}
		}(unsignedInfo)
	}
	wg.Wait() // 阻塞等待完成
	log.Println("[doSignEvening]", u.Name, "代签完成!")
	if u.QqBotRevue.Module == "brief" {
		u.qqBotRevueEvent("晚检代签提醒", "("+u.Name+")"+"代签人数为:", unsignedName)
	} else {
		u.qqBotRevueEvent("晚检代签提醒", "("+u.Name+")"+"代签名单为:", unsignedName)
	}
}

func (u User) EveningSignOperate() {
	signId := u.getEveningSignId()
	if signId == "0" {
		log.Println("[doSignEvening]", u.Name, "请求签到信息发生错误")
		u.qqBotRevueEvent("晚检代签提醒", "("+u.Name+")"+"请求签到信息发生错误")
	} else if signId == "1" {
		log.Println("[doSignEvening]", u.Name, "未到(或已过)签到时间")
		u.qqBotRevueEvent("晚检代签提醒", "("+u.Name+")"+"未到(或已过)签到时间")
	} else {
		unsignedList := u.getUnsignedList(signId)
		if len(unsignedList) != 0 {
			u.doSignEvening(unsignedList)
		} else {
			log.Println("[doSignEvening]", u.Name, "获取签到名单失败,可能所有同学已经签到")
			u.qqBotRevueEvent("晚检代签提醒", "("+u.Name+")"+"未到(或已过)签到时间")
		}
	}
}

func (u User) HealthCheckOperate() {
	log.Println("[HealthCheckOperate]", "代签完成!")
	client := resty.New()
	page := 1
	var unsignedStuId []string
	var unsignedName []string
	for {
		post, err := client.R().SetQueryParams(map[string]string{
			"date": getDate(),
			"page": strconv.Itoa(page),
			"size": "15",
		}).SetHeaders(map[string]string{
			"JWSESSION":  u.Jwsession,
			"User-Agent": u.UserAgent,
		}).Post("https://teacher.wozaixiaoyuan.com/health/getNoHealthUsers.json")
		if err != nil {
			log.Println("[HealthCheckOperate]", u.Name, "未打卡名单请求错误,错误信息为:"+err.Error())
			return
		}
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		if postInfo.Find("code") != -10 {
			if len(postInfo.Reset().Find("data").([]interface{})) == 0 {
				if page == 1 {
					log.Println("[HealthCheckOperate]seq=", u.Name, "没有打卡信息或者打卡没有开始!")
					//u.qqBotRevueEvent("日检日报代签提醒", "没有打卡信息或者打卡没有开始!")
					return
				}
				break
			}
			unsignedData := postInfo.Reset().From("data").Select("id", "name").Get()
			for _, data := range unsignedData.([]interface{}) {
				unsignedStuId = append(unsignedStuId, data.(map[string]interface{})["id"].(string))
				unsignedName = append(unsignedName, data.(map[string]interface{})["name"].(string))
			}
			page++
			//time.Sleep(1 * time.Second)
		} else {
			log.Println("[HealthCheckOperate]", u.Name, "jwsession失效,请更换!")
			u.qqBotRevueEvent("健康打卡代签提醒", "jwsession失效,请更换!")
			break
		}
	}
	time.Sleep(1 * time.Second)
	wg.Add(len(unsignedStuId))
	for i := 0; i < len(unsignedStuId); i++ {
		go u.healthCheckForStudent(unsignedName[i], unsignedStuId[i])
		//time.Sleep(1 * time.Second)
	}
	wg.Wait()

	log.Println("[HealthCheckOperate]", "打卡完成!")
	if u.QqBotRevue.Module == "brief" {
		u.qqBotRevueEvent("健康打卡代签提醒", "("+u.Name+")"+"代签人数为:", unsignedName)
	} else {
		u.qqBotRevueEvent("健康打卡代签提醒", "("+u.Name+")"+"代签名单为:", unsignedName)
	}

}

func (u User) healthCheckForStudent(stuName, stuId string) {
	defer wg.Done()

	province := "陕西省"
	city := "西安市"
	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒
	content := province + "_" + signTime + "_" + city
	signatureHeader := getSha256(content)
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetQueryParams(map[string]string{
		"answers":         "[\"0\"]",
		"temperature":     "36.5",
		"userId":          stuId,
		"latitude":        "",
		"longitude":       "",
		"country":         "中国",
		"city":            "西安市",
		"district":        "鄠邑区",
		"province":        "陕西省",
		"township":        "",
		"street":          "",
		"areacode":        "",
		"towncode":        "156",
		"timestampHeader": signTime,
		"signatureHeader": signatureHeader,
	}).Post("https://teacher.wozaixiaoyuan.com/health/save.json")
	if err != nil {
		log.Println("[healthCheckForStudent]", u.Name, stuName, "执行代打卡发生错误:"+err.Error())
		return
	}
	postJson := gojsonq.New().JSONString(string(post.Body()))
	if int(postJson.Reset().Find("code").(float64)) != 0 {
		log.Println("[healthCheckForStudent]", u.Name, stuName, "执行代打卡请求错误:"+string(post.Body()))
		return
	}
	//log.Println("[doSignEvening]", u.Name, stuName, "签到成功")
	//log.Println("[healthCheckForStudent]", stuName, "签到成功")

}
