package server

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	dataMutex = &sync.RWMutex{}
	dataFile  *os.File
	data      [][]string
)

// InitFile create data.csv if not exist
func InitFile() {
	_, err := os.Stat("data.csv")
	if err == nil {
		fmt.Println("file already exist")
	} else if os.IsNotExist(err) {
		// file not exist, create
		file, err := os.Create("data.csv")
		if err != nil {
			fmt.Println("crate file fail:", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		// write header
		header := []string{"C1", "C2", "C3"}
		err = writer.Write(header)
		if err != nil {
			fmt.Println("write header fail:", err)
			return
		}

		writer.Flush()

		fmt.Println("create file success")
	} else {
		fmt.Println("file exist judge fail:", err)
	}

	dataFile, err = os.OpenFile("data.csv", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer dataFile.Close()

	reader := csv.NewReader(dataFile)
	data, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	readCsv()
}

// test read csv after create
func readCsv() {
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println("open file fail:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("read file fail:", err)
		return
	}

	for _, row := range rows {
		fmt.Println(row)
	}
}

func getValue(p string) string {
	index := strings.LastIndex(p, "=")

	return p[index+1:]
}

func matchPredicate(row []string, predicate string) bool {
	// 解析谓词
	operator := ""
	term := ""
	if strings.Contains(predicate, "==") {
		operator = "=="
		// term = strings.Trim(predicate, "C1 C2 C3 == \"")
		term = getValue(predicate)
		fmt.Printf(">>>> test, in ==, predicate=: %s, term: %s\n", predicate, term)
	} else if strings.Contains(predicate, "!=") {
		operator = "!="
		term = getValue(predicate)
	} else if strings.Contains(predicate, "$=") {
		operator = "$="
		term = getValue(predicate)
	} else if strings.Contains(predicate, "&=") {
		operator = "&="
		term = getValue(predicate)
	}

	// 检查谓词是否匹配
	switch operator {
	case "==":
		return rowMatches(row, term)
	case "!=":
		return !rowMatches(row, term)
	case "$=":
		return rowMatchesIgnoreCase(row, term)
	case "&=":
		return rowContains(row, term)
	default:
		return false
	}
}

func rowMatches(row []string, term string) bool {
	for _, cell := range row {
		if cell == term {
			return true
		}
	}
	return false
}

func rowMatchesIgnoreCase(row []string, term string) bool {
	for _, cell := range row {
		if strings.EqualFold(cell, term) {
			return true
		}
	}
	return false
}

func rowContains(row []string, term string) bool {
	for _, cell := range row {
		if strings.Contains(cell, term) {
			return true
		}
	}
	return false
}

func executeUpdate(values []string) error {
	//if len(values) != 6 {
	//	return fmt.Errorf("invalid UPDATE command")
	//}

	fmt.Printf("update values: %v\n", values)
	var newData [][]string

	for _, row := range data {
		if rowMatchesUpdate(row, values) {
			row = updateRow(row, values)
		}
		newData = append(newData, row)
	}

	data = newData

	// 更新数据文件
	err := updateDataFile()
	if err != nil {
		return err
	}

	return nil
}

func rowMatchesUpdate(row []string, values []string) bool {
	for i := 0; i < len(values)-2; i++ {
		if values[i] != "*" && row[i] != values[i] {
			return false
		}
	}
	return true
}

func updateRow(row []string, values []string) []string {
	for i := 3; i < 6; i += 2 {
		if values[i] != "*" {
			row[i/2-1] = values[i]
		}
	}
	return row
}

func getCommand(p string) (command string, values []string) {
	// fmt.Printf("job modify, p: %s\n", p)
	index := strings.Index(p, " ")
	command = p[:index]
	values = strings.Split(p[index+1:], ",")
	return
}

func executeModify(job string) error {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	// 解析修改命令
	cmd, vals := getCommand(job)
	// fmt.Printf("vals: %v\n", vals)
	//command := strings.Split(job, " ")
	//val := strings.Split(command[1], ",")

	switch cmd {
	case "INSERT":
		return executeInsert(vals)
	case "DELETE":
		return executeDelete(vals)
	case "UPDATE":
		return executeUpdate(vals)
	default:
		return fmt.Errorf("invalid modification command")
	}
}

func matchPredicates(row []string, predicates []string) bool {
	for _, predicate := range predicates {
		if !matchPredicate(row, predicate) {
			return false
		}
	}
	return true
}

func executeInsert(values []string) error {
	if len(values) != 3 {
		fmt.Printf("values: %v, len: %d\n", values, len(values))
		return fmt.Errorf("invalid INSERT command")
	}

	data = append(data, values)

	// 更新数据文件
	err := updateDataFile()
	if err != nil {
		return err
	}

	return nil
}

func executeDelete(values []string) error {
	if len(values) < 1 || len(values) > 3 {
		return fmt.Errorf("invalid DELETE command")
	}

	var newData [][]string
	for _, row := range data {
		if !rowMatchesDelete(row, values) {
			newData = append(newData, row)
		}
	}

	data = newData

	// 更新数据文件
	err := updateDataFile()
	if err != nil {
		return err
	}

	return nil
}

func rowMatchesDelete(row []string, values []string) bool {
	for i, value := range values {
		if row[i] != value {
			return false
		}
	}
	return true
}

func updateDataFile() error {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "temp_data.csv")
	if err != nil {
		return err
	}
	// defer tempFile.Close()

	// 写入数据到临时文件
	writer := csv.NewWriter(tempFile)
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	// 关闭临时文件
	writer.Flush()
	tempFile.Close()

	err = copyFile(tempFile.Name(), dataFile.Name())
	if err != nil {
		return err
	}

	// 删除临时文件
	err = os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func executeQuery(query string) [][]string {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	// 解析查询条件
	predicates := strings.Split(query, " and ")

	fmt.Printf("predicates: %v\n", predicates)
	// 执行查询操作
	var result [][]string
	for _, row := range data {
		if matchPredicates(row, predicates) {
			result = append(result, row)
		}
	}

	return result
}
