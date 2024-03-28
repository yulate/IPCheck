/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"IPCheck/cmd"
	"IPCheck/utils/iploc"
	"embed"
	"fmt"
	"os"
)

const (
	filePath = "data/czutf8.dat"
)

func main() {
	cmd.Execute()
}

func init() {
	CheckData()
	iploc.IpLocInit()
}

func CheckData() {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("文件不存在, 释放中...")
		err := os.Mkdir("data", 777)
		if err != nil {
			fmt.Println("创建data文件夹失败：", err)
		}
		ReleaseTheFile()
	}
}

//go:embed data/czutf8.dat
var data embed.FS

func ReleaseTheFile() {
	file, err := data.ReadFile("data/czutf8.dat")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filePath, []byte(file), 0644)
	if err != nil {
		fmt.Println("Failed to write file:", err)
		return
	}
}
