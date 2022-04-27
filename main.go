package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
	"time"
)

var (
	cmd = tea.NewProgram(newModel()) //定义全局变量
)

func pause() {
	flag := 5
	for i := 0; i < 5; i++ {
		cmd.Send(pauseInfo{info: "程序将在" + strconv.Itoa(flag-i) + "秒后自动关闭...\n"})
		time.Sleep(1 * time.Second)
	}
	cmd.Send(tea.Quit())
}
func main() {
	go func() {
		if configFileExists, _ := pathExists("./config.json"); !configFileExists {
			cmd.Send(printInfo{code: 2, funcName: "main", info: "没有找到config.json文件,请将config.json放入该程序同一文件夹下！"})
			pause()
		}
	}()
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}
