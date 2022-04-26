package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type printInfo struct {
	funcName string
	name     string
	status   string
	code     int // 如果code==0,则显示以上信息
	info     string
}

var (
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	cmd          = tea.NewProgram(newModel()) //定义全局变量
)

func (r printInfo) String() (s string) {
	if r.code == 1 {
		s += fmt.Sprintf("function: %s | name: %s | status: %s \n", r.funcName, r.name, r.status)
		s += fmt.Sprintf("info: %s\n", r.info)
	} else if r.code == 2 {
		s += fmt.Sprintf("function: %s | info: %s\n", r.funcName, r.info)
	} else if r.code == 3 {
		s += fmt.Sprintf("function: %s | name: %s | status: %s \n", r.funcName, r.name, r.status)
	}
	return s
}

func newModel() model {
	const numLastResults = 5 // 最大消息列表显示数量
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		items:         []string{"晨检打卡", "午检打卡", "晚检签到"},
		index:         0,
		doTask:        false,
		printInfoList: make([]printInfo, numLastResults),
		spinner:       s,
	}

}

func main() {
	if configFileExists, _ := pathExists("./config.json"); !configFileExists {
		fmt.Println("没有找到config.json文件,请将config.json放入该程序同一文件夹下！")
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}

}
