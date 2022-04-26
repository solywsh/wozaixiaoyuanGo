package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

//需要实现init update view三个方法
type model struct {
	items         []string
	index         int
	spinner       spinner.Model
	printInfoList []printInfo
	doTask        bool
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
				go func() { dailyCheck(1) }()
			} else if m.index == 1 {
				go func() { dailyCheck(2) }()
			} else {
				go func() { eveningSignOperate() }()
			}
			return m, nil
			//return m, tea.Quit
		}
	case printInfo:
		m.printInfoList = append(m.printInfoList[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, spinnerCmd
	}
	return m, nil
}

//渲染列表
func (m model) View() string {
	s := "我在校园打卡(导员版)\n\n"

	if m.doTask {
		s += m.spinner.View() + " 正在执行中...\n"
	} else {
		s += "请选择:\n\n"
	}

	for i, item := range m.items {
		selected := " "
		if m.index == i {
			selected = "»"
		}
		s += fmt.Sprintf("%s %s\n", selected, item)
	}

	if len(m.printInfoList) > 0 {
		for _, res := range m.printInfoList {
			if res.code != 0 {
				s += res.String()
			}
		}
	}
	s += "\n按住 Ctrl+C 或 Q 退出\n"
	return appStyle.Render(s)
}
