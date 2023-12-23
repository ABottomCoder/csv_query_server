package repository

import "testing"

var filePath = "../../pkg/csv/data.csv"

func TestInitFileInit(t *testing.T) {
	InitFile(filePath)
}
