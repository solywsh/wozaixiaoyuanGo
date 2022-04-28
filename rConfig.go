package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/thedevsaddam/gojsonq"
)

type config struct {
	jwsession string
	userAgent string
	province  string
	city      string
}

type yamlDecoder struct {
}

var yamlConfig config
var yamlPath = "./config.yaml"

func readConfig() (res string) {
	strHeader := "在" + yamlPath + "中没有找到"
	strFoot := ",请重新配置!"
	jq := gojsonq.New(gojsonq.SetDecoder(&yamlDecoder{})).File(yamlPath)
	if jq.Reset().Find("jwsession") == nil {
		return strHeader + "jwsession或为空" + strFoot
	} else {
		yamlConfig.jwsession = jq.Reset().Find("jwsession").(string)
	}
	if jq.Reset().Find("userAgent") == nil {
		return strHeader + "userAgent或为空" + strFoot
	} else {
		yamlConfig.userAgent = jq.Reset().Find("userAgent").(string)
	}
	if jq.Reset().Find("province") == nil {
		return strHeader + "province或为空" + strFoot
	} else {
		yamlConfig.province = jq.Reset().Find("province").(string)
	}
	if jq.Reset().Find("city") == nil {
		return strHeader + "city或为空" + strFoot
	} else {
		yamlConfig.city = jq.Reset().Find("city").(string)
	}
	return "0"
}

// Decode 实现gojsonq.Decoder
func (i *yamlDecoder) Decode(data []byte, v interface{}) error {
	bb, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb, &v)
}

//func main() {
//	fmt.Println(readConfig())
//}
