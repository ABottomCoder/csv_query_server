package utils

import (
	"fmt"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	path := GetFilePath()

	fmt.Printf("path: %s\n", path)
}
