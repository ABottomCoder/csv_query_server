package utils

import (
	"fmt"
	"os"
)

func GetFilePath() string {
	if len(os.Args) < 3 {
		fmt.Println("请提供一个路径作为命令行参数")
		return ""
	}

	path := os.Args[2]

	return path
}
