package main

import (
	"crypto/sha256"
	"encoding/hex"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func getDate() string {
	// 表单提交使用
	return time.Now().Format("20060102")
}

func getDataAndTime() string {
	// 验证时间使用
	return time.Now().Format("2006-01-02 15:04")
}

func getSha256(src string) string {
	sha256Bytes := sha256.Sum256([]byte(src))
	sha256String := hex.EncodeToString(sha256Bytes[:])
	return sha256String
}

func checkForStudent(stuName, stuId, seq, jwsession, userAgent string) {

	province := yamlConfig.province
	city := yamlConfig.city
	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒

	content := province + "_" + signTime + "_" + city
	signatureHeader := getSha256(content)
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  jwsession,
		"User-Agent": userAgent,
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
		cmd.Send(printInfo{code: 2, funcName: "checkForStudent", info: "执行代打卡发生错误:" + err.Error()})
		return
	}
	msg := printInfo{
		funcName: "checkForStudent",
		name:     stuName,
	}
	postJson := gojsonq.New().JSONString(string(post.Body()))
	if int(postJson.Reset().Find("code").(float64)) == 0 {
		msg.code = 1
		msg.status = "正常"
	} else {
		msg.code = 3
		msg.status = "失败"
		msg.info = "发生错误，错误信息为：" + postJson.Reset().Find("message").(string)
	}
	cmd.Send(msg)
	//fmt.Println(stuId, string(post.Body()))
	wg.Done()
}

func dailyCheck(seq int) {
	date := getDate()
	jwsession := yamlConfig.jwsession
	userAgent := yamlConfig.userAgent

	client := resty.New()
	page := 1
	var unsignedStuId []string
	var unsignedName []string
	for {
		post, err := client.R().SetQueryParams(map[string]string{
			"seq":  strconv.Itoa(seq),
			"date": date,
			"type": "0",
			"page": strconv.Itoa(page),
			"size": "20",
		}).SetHeaders(map[string]string{
			"JWSESSION":  jwsession,
			"User-Agent": userAgent,
		}).Post("https://teacher.wozaixiaoyuan.com/heat/getHeatUsers.json")
		if err != nil {
			cmd.Send(printInfo{code: 2, funcName: "dailyCheck", info: "未打卡名单请求错误,错误信息为:" + err.Error()})
			return
		}
		//postMap := JsonByteToMap(post.Body())
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		if postInfo.Find("code") != -10 {
			if len(postInfo.Reset().Find("data").([]interface{})) == 0 {
				if page == 1 {
					cmd.Send(printInfo{code: 2, funcName: "dailyCheck", info: "没有打卡信息或者打卡没有开始!"})
					pause()
					cmd.Send(tea.Quit())
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
			cmd.Send(printInfo{code: 2, funcName: "dailyCheck", info: "jwsession失效,请更换!"})
			pause()
			cmd.Send(tea.Quit())
			break
		}
	}
	cmd.Send(printInfo{code: 2, funcName: "dailyCheck", info: "开始执行打卡.."})
	time.Sleep(1 * time.Second)
	wg.Add(len(unsignedStuId))
	for i := 0; i < len(unsignedStuId); i++ {
		go checkForStudent(unsignedName[i], unsignedStuId[i], strconv.Itoa(seq), jwsession, userAgent)
		//time.Sleep(1 * time.Second)
	}
	wg.Wait()
	cmd.Send(printInfo{code: 2, funcName: "dailyCheck", info: "签到完成!"})
	pause()
	cmd.Send(tea.Quit())
}

func getEveningSignId(jwsession string) (signId string) {
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION": jwsession,
	}).SetQueryParams(map[string]string{
		"state":   "0",
		"keyword": "",
		"page":    "1",
	}).Post("https://teacher.wozaixiaoyuan.com/signManage/getList.json")
	if err != nil {
		cmd.Send(printInfo{code: 2, funcName: "getEveningSignId", info: "请求晚检签到id发生错误,错误信息为:" + err.Error()})
		return "0" //发生错误
	}
	//postInfo := JsonByteToMap(post.Body())
	//fmt.Println(string(post.Body()))
	postInfo := gojsonq.New().JSONString(string(post.Body()))
	if int(postInfo.Find("code").(float64)) == 0 {
		signEndTime := postInfo.Reset().Find("data.[0].end").(string)
		signBeginTime := postInfo.Reset().Find("data.[0].start").(string)
		signId := postInfo.Reset().Find("data.[0].id").(string)
		signNowTime := getDataAndTime()
		if signNowTime > signBeginTime && signEndTime > signNowTime {
			//fmt.Println(signId)
			return signId // 正常返回签到Id
		} else {
			return "1" // 未到(或已过)签到时间
		}
	}
	return "0"
}

func getUnsignedList(signId, jwsession string) (unsignedList []map[string]interface{}) {
	url := "https://teacher.wozaixiaoyuan.com/signManage/getSignLogs.json"
	client := resty.New()
	page := 1
	for {
		post, err := client.R().SetHeaders(map[string]string{
			"JWSESSION": jwsession,
		}).SetQueryParams(map[string]string{
			"id":       signId,
			"page":     strconv.Itoa(page),
			"type":     "0",
			"size":     "20",
			"targetId": "",
		}).Post(url)
		if err != nil {
			cmd.Send(printInfo{code: 2, funcName: "getUnsignedList", info: "请求晚检未签到名单发生错误,错误信息为:" + err.Error()})
			pause()
			cmd.Send(tea.Quit())
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

func doSignEvening(unsignedList []map[string]interface{}, jwsession string) {
	url := "https://teacher.wozaixiaoyuan.com/signManage/adminSign.json"
	for _, unsignedInfo := range unsignedList {
		client := resty.New()
		post, err := client.R().SetHeaders(map[string]string{
			"JWSESSION": jwsession,
		}).SetQueryParams(map[string]string{
			"id":   unsignedInfo["id"].(string),
			"type": "1",
		}).Post(url)
		if err != nil {
			msg := printInfo{
				code:     3,
				funcName: "doSignEvening",
				name:     unsignedInfo["name"].(string),
				info:     "执行晚检代签发生错误,错误信息为:" + err.Error(),
				status:   "错误",
			}
			cmd.Send(msg)
			return
		}
		msg := printInfo{
			funcName: "doSignEvening",
			name:     unsignedInfo["name"].(string),
		}
		rJson := gojsonq.New().JSONString(string(post.Body()))
		if int(rJson.Reset().Find("code").(float64)) == 0 {
			msg.code = 1
			msg.status = "正常"
		} else {
			msg.code = 3
			msg.status = "失败"
			msg.info = "代签失败，失败信息为" + rJson.Reset().Find("message").(string)
		}
		cmd.Send(msg)
		wg.Done()
	}
}

func eveningSignOperate() {
	jwsession := yamlConfig.jwsession
	signId := getEveningSignId(jwsession)
	if signId == "0" {
		cmd.Send(printInfo{code: 2, funcName: "eveningSignOperate", info: "请求签到信息发生错误"})
		pause()
		cmd.Send(tea.Quit())
	} else if signId == "1" {
		cmd.Send(printInfo{code: 2, funcName: "eveningSignOperate", info: "未到(或已过)签到时间"})
		pause()
		cmd.Send(tea.Quit())
	} else {
		unsignedList := getUnsignedList(signId, jwsession)
		if len(unsignedList) != 0 {
			// 执行签到
			cmd.Send(printInfo{code: 2, funcName: "eveningSignOperate", info: "开始代签..."})
			time.Sleep(1 * time.Second)
			wg.Add(len(unsignedList))
			go doSignEvening(unsignedList, jwsession)
			wg.Wait() // 阻塞等待完成
			cmd.Send(printInfo{code: 2, funcName: "eveningSignOperate", info: "代签完成!"})
			pause()
			cmd.Send(tea.Quit())
		} else {
			cmd.Send(printInfo{code: 2, funcName: "eveningSignOperate", info: "获取签到名单失败,可能所有同学已经签到"})
			pause()
			cmd.Send(tea.Quit())
		}
	}
}
