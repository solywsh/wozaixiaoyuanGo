package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"os"
	"strconv"
	"time"
)

func sendMessagePlusToken(title, content string) {
	client := resty.New()
	configInfo := gojsonq.New().File("./config.json")
	pushPlusTokenList := configInfo.Reset().Find("pushPlusToken").([]interface{})
	for _, token := range pushPlusTokenList {
		post, err := client.R().SetHeader("Content-Type", "application/json").
			SetQueryParams(map[string]string{
				"token":   token.(string),
				"title":   title,
				"content": content,
			}).Post("https://pushplus.hxtrip.com/send")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(post.Body()))
	}
}

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

func checkForStudent(stuId, seq, jwsession, userAgent string) {

	configInfo := gojsonq.New().File("./config.json")
	province := configInfo.Reset().Find("province").(string)
	city := configInfo.Reset().Find("city").(string)

	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒

	content := province + "_" + signTime + "_" + city
	signatureHeader := getSha256(content)
	//postData = "answers=%5B%220%22%5D&seq=" +
	//	strconv.Itoa(seq) +
	//	"&temperature=36.0&userId=" +
	//	strconv.Itoa(stuId) +
	//	"&latitude=&longitude=&country=&city=&district=&province=&township=&street=&areacode=&timestampHeader=" +
	//	signTime +
	//	"&signatureHeader=" +
	//	signatureHeader
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
		panic(err)
	}
	fmt.Println(stuId, string(post.Body()))
}

func dailyCheck(seq int) {
	date := getDate()
	configInfo := gojsonq.New().File("./config.json")
	jwsession := configInfo.Reset().Find("jwsession").(string)
	userAgent := configInfo.Reset().Find("userAgent").(string)

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
			panic(err)
		}
		//postMap := JsonByteToMap(post.Body())
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		if postInfo.Find("code") != -10 {
			if len(postInfo.Reset().Find("data").([]interface{})) == 0 {
				if page == 1 {
					fmt.Println("没有打卡信息或者打卡没有开始!")
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
			fmt.Println("jwsession失效")
			break
		}
	}
	fmt.Println("开始执行打卡...")
	for i := 0; i < len(unsignedStuId); i++ {
		fmt.Println(unsignedName[i], unsignedStuId[i])
		checkForStudent(unsignedStuId[i], strconv.Itoa(seq), jwsession, userAgent)
		//time.Sleep(1 * time.Second)
	}
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
		fmt.Println(err)
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
			panic(err)
		}
		postInfo := gojsonq.New().JSONString(string(post.Body()))
		//fmt.Println(reflect.TypeOf(postInfo.Find("data.[0]")))
		if int(postInfo.Reset().Find("code").(float64)) == 0 && len(postInfo.Reset().Find("data").([]interface{})) != 0 {
			selectInfo := postInfo.Reset().From("data").Select("name", "id").Get() // 只提取出id和name
			for _, info := range selectInfo.([]interface{}) {
				unsignedList = append(unsignedList, info.(map[string]interface{}))
				fmt.Println(info) // map[id:437356871436734465 name:马鹏]
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
			panic(err)
		}
		fmt.Println(unsignedInfo["name"], string(post.Body()))
	}
}

func eveningSignOperate() {
	// 读取配置文件
	configInfo := gojsonq.New().File("./config.json")
	jwsession := configInfo.Reset().Find("jwsession").(string)
	signId := getEveningSignId(jwsession)
	var unsignedName string //未打卡的姓名列表
	//signId := "437356864029593600"
	if signId == "0" {
		fmt.Println("请求签到信息发生错误")
	} else if signId == "1" {
		fmt.Println("未到(或已过)签到时间")
	} else {
		unsignedList := getUnsignedList(signId, jwsession)
		if len(unsignedList) != 0 {
			// 执行签到
			doSignEvening(unsignedList, jwsession)
			for _, info := range unsignedList {
				unsignedName += info["name"].(string) + " "
			}
			fmt.Println(unsignedName)
		} else {
			fmt.Println("获取签到名单失败,可能所有同学已经签到")
		}
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true, nil
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//其它类型，不确定是否存在
	return false, err
}

func main() {
	if configFileExists, _ := pathExists("./config.json"); !configFileExists {
		fmt.Println("没有找到config.json文件,请将config.json放入该程序同一文件夹下！")
		return
	}
	//dailyCheck("20220422", 1)
	//sendMessagePlusToken("测试标题", "测试内容2")
	//getEveningSignList("afbfcb940699498bb76eb5389ac483fd")
	//getUnsignedList("437356864029593600", "afbfcb940699498bb76eb5389ac483fd")
	//eveningSignOperate("afbfcb940699498bb76eb5389ac483fd")
}
