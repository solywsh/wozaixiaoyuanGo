package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type printInfo struct {
	funcName string
	name     string
	status   string
	code     int // 如果code==0,则显示以上信息
	info     string
}

type pauseInfo struct {
	info string
}

var (
	appStyle        = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	spinnerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))      //定义全局变量
	crimson         = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC143C")) //暗红 error
	lawngreen       = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CFC00")) //草绿色 normal
	yellow          = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")) //黄色 warring/info
	mediumblue      = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000CD")) //间兰色 func
	mediumturquoise = lipgloss.NewStyle().Foreground(lipgloss.Color("#48D1CC")) //亮绿色 keyword
)

func (r printInfo) String() (s string) {
	if r.code == 1 {
		var statusString string
		if r.status == "正常" {
			statusString = lawngreen.Render(r.status)
		} else {
			statusString = crimson.Render(r.status)
		}
		s += fmt.Sprintf("%s: %s | %s: %s | %s: %s \n",
			mediumturquoise.Render("function"), mediumblue.Render(r.funcName),
			mediumturquoise.Render("name"), r.name,
			mediumturquoise.Render("status"), statusString)
	} else if r.code == 2 {
		s += fmt.Sprintf("%s: %s | %s: %s\n",
			mediumturquoise.Render("function"), mediumblue.Render(r.funcName),
			mediumturquoise.Render("info"), yellow.Render(r.info))
	} else if r.code == 3 {
		var statusString string
		if r.status == "正常" {
			statusString = lawngreen.Render(r.status)
		} else {
			statusString = crimson.Render(r.status)
		}
		s += fmt.Sprintf("%s: %s | %s: %s | %s: %s \n",
			mediumturquoise.Render("function"), r.funcName,
			mediumturquoise.Render("name"), r.name,
			mediumturquoise.Render("status"), statusString)
		s += fmt.Sprintf("%s: %s\n",
			mediumturquoise.Render("info"), yellow.Render(r.info))
	}
	return s
}

//需要实现init update view三个方法
type model struct {
	items         []string
	index         int
	spinner       spinner.Model //加载动画
	printInfoList []printInfo
	doTask        bool      // 控制是否显示加载动画
	pause         pauseInfo //暂停时显示的信息
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		//当输入q或者ctrl+c 退出
		case "ctrl+c", "q":
			return m, tea.Quit
		// 如果是up 光标向上移动
		case "up":
			if m.index > 0 {
				m.index--
			}
		//如果是down 光标向下移动
		case "down":
			if m.index < len(m.items)-1 {
				m.index++
			}
		//如果是enter 处理事件
		case "enter", " ":
			//fmt.Println(m.items[m.index])
			m.doTask = true
			if m.index == 0 {
				m.index = 4
				// 使用异步的方式去启动
				go dailyCheck(1)
			} else if m.index == 1 {
				m.index = 4
				go dailyCheck(2)
			} else {
				m.index = 4
				go eveningSignOperate()
			}
			return m, nil
			//return m, tea.Quit
		}
	case printInfo:
		// 接收打印信号
		m.printInfoList = append(m.printInfoList[1:], msg)
		return m, nil
	case spinner.TickMsg:
		// 接收进度动画信号
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, spinnerCmd
	case pauseInfo:
		// 接收暂停显示信号
		m.doTask = true
		m.pause = msg
		return m, nil
	}
	return m, nil
}

// View 渲染列表
func (m model) View() string {
	s := "WoZaiXiaoYuanGo For Teacher\n\n"
	if m.doTask {
		s += m.spinner.View() + " 正在执行中...\n"
	} else {
		s += "请选择:\n"
		for i, item := range m.items {
			selected := " "
			if m.index == i {
				selected = "»"
			}
			s += fmt.Sprintf("%s %s\n", selected, item)
		}
	}
	if len(m.printInfoList) > 0 {
		for _, res := range m.printInfoList {
			if res.code != 0 {
				s += res.String()
			}
		}
	}
	if m.pause.info != "" {
		s += yellow.Render(m.pause.info)
	}
	s += "\n按住 Ctrl+C 或 Q 退出\n"
	s += "power by https://github.com/solywsh/wozaixiaoyuanGo"
	return appStyle.Render(s)
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
