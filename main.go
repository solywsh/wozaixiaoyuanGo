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
	flag := 4
	for i := 0; i < 4; i++ {
		cmd.Send(pauseInfo{info: "程序将在" + strconv.Itoa(flag-i) + "秒后自动关闭..."})
		time.Sleep(1 * time.Second)
	}
}
func main() {
	if configFileExists, _ := pathExists("./config.json"); !configFileExists {
		fmt.Println("没有找到config.json文件,请将config.json放入该程序同一文件夹下！")
		pause()
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}

}
