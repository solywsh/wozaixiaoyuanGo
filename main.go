package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

var cmd = tea.NewProgram(newModel()) //定义全局变量

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
	go func() {
		if configFileExists, _ := pathExists("./config.yaml"); !configFileExists {
			cmd.Send(printInfo{code: 2, funcName: "main", info: "没有找到config.yaml文件,请将config.yaml放入该程序同一文件夹下！"})
			pause()
		}
		if res := readConfig(); res != "0" {
			cmd.Send(printInfo{code: 2, funcName: "main", info: res})
			pause()
		}
	}()
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}
