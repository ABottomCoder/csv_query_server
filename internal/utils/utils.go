package utils

import (
	"fmt"
	"os"
)

func GetFilePath() string {
	if len(os.Args) < 2 {
		fmt.Println("请提供一个路径作为命令行参数")
		return ""
	}

	path := os.Args[1]

	return path
}
