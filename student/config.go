package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)

type Config struct {
	User []User `yaml:"user"`
}

type EveningCheck struct {
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
}

type User struct {
	Jwsession      string         `yaml:"jwsession"`
	QqBotRevue     QqBotRevue     `yaml:"qqBotRevue"`
	Name           string         `yaml:"name"`
	MorningCheck   MorningCheck   `yaml:"morningCheck"`
	AfternoonCheck AfternoonCheck `yaml:"afternoonCheck"`
	EveningCheck   EveningCheck   `yaml:"eveningCheck"`
	Province       string         `yaml:"province"`
	City           string         `yaml:"city"`
	UserAgent      string         `yaml:"userAgent"`
}

type MorningCheck struct {
	EndTime   string `yaml:"endTime"`
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
}

type AfternoonCheck struct {
	Enable    bool   `yaml:"enable"`
	CheckTime string `yaml:"checkTime"`
	EndTime   string `yaml:"endTime"`
}

type QqBotRevue struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
	UserId string `yaml:"userId"`
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
