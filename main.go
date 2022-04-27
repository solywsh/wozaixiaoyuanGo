package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

var (
	cmd = tea.NewProgram(newModel()) //定义全局变量
)

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
