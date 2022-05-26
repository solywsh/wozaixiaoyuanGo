package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)

type Config struct {
	User []User `yaml:"user"`
}

type AfternoonCheck struct {
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
}

type User struct {
	Jwsession      string         `yaml:"jwsession"`
	Province       string         `yaml:"province"`
	City           string         `yaml:"city"`
	Name           string         `yaml:"name"`
	UserAgent      string         `yaml:"userAgent"`
	MorningCheck   MorningCheck   `yaml:"morningCheck"`
	AfternoonCheck AfternoonCheck `yaml:"afternoonCheck"`
	EveningCheck   EveningCheck   `yaml:"eveningCheck"`
	HealthCheck    HealthCheck    `yaml:"healthCheck"`
	QqBotRevue     QqBotRevue     `yaml:"qqBotRevue"`
}

type QqBotRevue struct {
	Enable bool   `yaml:"enable"`
	Module string `yaml:"module"`
	Qq     []Qq   `yaml:"qq"`
}

type Qq struct {
	UserId string `yaml:"userId"`
	Token  string `yaml:"token"`
}

type HealthCheck struct {
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
}

type EveningCheck struct {
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
	Enable    bool   `yaml:"enable"`
}

type MorningCheck struct {
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
}

func NewConf(yamlPath string) (conf Config, err error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Println("文件打开错误,请传入正确的文件路径!", err)
		return conf, err
	}
	//fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &conf)
	//err = yaml.UnmarshalStrict(yamlFile, kafkaCluster)
	if err != nil {
		log.Println("文件解析错误,请配置正确的yaml格式!", err)
		return conf, err
	}
	return conf, nil
}

func main() {
	c, _ := NewConf("./config.yaml")
	for _, v := range c.User {
		for _, qq := range v.QqBotRevue.Qq {
			fmt.Println(v.Name, qq)
		}
	}
}
