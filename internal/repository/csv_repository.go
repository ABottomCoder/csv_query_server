package repository

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
	dataMutex       = &sync.RWMutex{}
	dataFile        *os.File
	data            [][]string
	Headers         []string
	Header2Index    = make(map[string]int)
	DefaultFilePath = "pkg/csv/data.csv"
)

// InitFile create data.csv if not exist
func InitFile(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		fmt.Println("file already exist")

	} else if os.IsNotExist(err) {
		// file not exist, create
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("crate file fail:", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		// write header
		defaultHeaders := []string{"C1", "C2", "C3"}
		err = writer.Write(defaultHeaders)
		if err != nil {
			fmt.Println("write header fail:", err)
			return
		}

		writer.Flush()

		fmt.Println("create file success")
	} else {
		fmt.Println("file exist judge fail:", err)
	}

	dataFile, err = os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer dataFile.Close()

	reader := csv.NewReader(dataFile)
	data, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	Headers = data[0]
	for i, header := range Headers {
		Header2Index[header] = i
	}

	fmt.Printf("Headers: %v\n", Headers)

	fmt.Printf("data: %v\n", data)

	// readCsv()
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

func matchPredicate(row []string, colIndex int, query [2]string) bool {
	// 解析谓词
	operator := query[0]
	term := query[1]
	//if strings.Contains(predicate, "==") {
	//	operator = "=="
	//	// term = strings.Trim(predicate, "C1 C2 C3 == \"")
	//	term = getValue(predicate)
	//	fmt.Printf(">>>> test, in ==, predicate=: %s, term: %s\n", predicate, term)
	//} else if strings.Contains(predicate, "!=") {
	//	operator = "!="
	//	term = getValue(predicate)
	//} else if strings.Contains(predicate, "$=") {
	//	operator = "$="
	//	term = getValue(predicate)
	//} else if strings.Contains(predicate, "&=") {
	//	operator = "&="
	//	term = getValue(predicate)
	//}

	// 检查谓词是否匹配
	switch operator {
	case "==":
		return rowMatches(row, colIndex, term)
	case "!=":
		return !rowMatches(row, colIndex, term)
	case "$=":
		return rowMatchesIgnoreCase(row, colIndex, term)
	case "&=":
		return rowContains(row, colIndex, term)
	default:
		return false
	}
}

func rowMatches(row []string, colIndex int, term string) bool {
	return row[colIndex] == term
}

func rowMatchesIgnoreCase(row []string, colIndex int, term string) bool {
	if strings.EqualFold(row[colIndex], term) {
		return true
	}

	return false
}

func rowContains(row []string, colIndex int, term string) bool {
	if strings.Contains(row[colIndex], term) {
		return true
	}

	return false
}

func getUpdateData(values []string) (idx int, newVal string) {
	idx = -1
	for i, header := range Headers {
		if header == values[len(values)-2] {
			idx = i
			break
		}
	}

	newVal = values[len(values)-1]

	return
}

func executeUpdate(values []string) error {
	idx, newVal := getUpdateData(values)
	if idx == -1 {
		return fmt.Errorf("invalid UPDATE command")
	}

	fmt.Printf("update values: %v\n", values)

	var newData [][]string

	for _, row := range data {
		if rowMatchesUpdate(row, values) {
			row = updateRow(row, idx, newVal)
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

func updateRow(row []string, idx int, newVal string) []string {
	row[idx] = newVal
	return row
}

func getCommand(p string) (command string, values []string) {
	// fmt.Printf("job modify, p: %s\n", p)
	index := strings.Index(p, " ")
	command = p[:index]
	values = strings.Split(p[index+1:], ",")
	return
}

func ExecuteModify(job string) error {
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

func matchPredicates(row []string, queryData map[int][2]string) bool {
	for conIndex, query := range queryData {
		if !matchPredicate(row, conIndex, query) {
			return false
		}
	}
	return true
}

func executeInsert(values []string) error {
	if len(values) != len(Headers) {
		fmt.Printf("values: %v, len: %d\n", values, len(values))
		return fmt.Errorf("invalid INSERT command")
	}

	// 先查询数据是否已存在
	queryData := make(map[int][2]string, len(Headers))
	for idx, _ := range Headers {
		queryData[idx] = [2]string{"==", values[idx]}
	}
	var result [][]string
	for _, row := range data {
		if matchPredicates(row, queryData) {
			result = append(result, row)
		}
	}

	if len(result) != 0 {
		return fmt.Errorf("data already exist")
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

func ExecuteQuery(query string) [][]string {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	// 解析查询条件
	predicates := strings.Split(query, " and ")

	fmt.Printf("predicates: %v\n", predicates)
	queryData := getQueryData(predicates)
	fmt.Printf("queryData: %v\n", queryData)

	// 执行查询操作
	var result [][]string
	for _, row := range data {
		if matchPredicates(row, queryData) {
			result = append(result, row)
		}
	}

	return result
}

// map[int][2]string ==> colindex->[operator, value]
func getQueryData(ps []string) (queryData map[int][2]string) {
	queryData = make(map[int][2]string, len(ps))
	for _, p := range ps {
		index := strings.LastIndex(p, "=")

		colName := p[:index-1]
		colValue := p[index+1:]
		if colIndex, ok := Header2Index[colName]; ok {
			queryData[colIndex] = [2]string{p[index-1 : index+1], colValue}
		}
	}

	return
}
