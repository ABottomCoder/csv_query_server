package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func RegisterModifyHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", modifyHandler)
}

type ModifyResponse struct {
	Msg string `json:"msg"`
}

// http://localhost:7259/?job=INSERT%20a1%2Ca2%2Ca3
// http://localhost:7259/?job=DELETE%20a1
// http://localhost:7259/?job=DELETE%20a1%2Ca2
// http://localhost:7259/?job=UPDATE%20ab%2Cac%2CC2%2Caa
// http://localhost:7259/?job=UPDATE%20a1%2Ca3%2CC3%2Ctt
func modifyHandler(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")

	// 检查输入
	if isValidModify(job) {
		response := ModifyResponse{}
		err := executeModify(job)
		if err != nil {
			response.Msg = "Execute modify fail"
			fmt.Printf("Execute modify fail, err: %v\n", err)
			jsonResponse(w, http.StatusBadRequest, response)
		} else {
			response.Msg = "Modification successful"
			jsonResponse(w, http.StatusOK, response)
		}
	} else {
		response := ModifyResponse{
			Msg: "Error description of the malformed modify",
		}
		jsonResponse(w, http.StatusBadRequest, response)
	}
}

func isValidModify(statement string) (match bool) {
	// 定义有效语句的正则表达式
	validStatementRegex := `^(DELETE|UPDATE|INSERT).*`

	// 使用正则表达式匹配输入语句
	match, _ = regexp.MatchString(validStatementRegex, statement)
	if !match {
		return
	}

	// 检查INSERT语句中是否包含所有列
	if strings.HasPrefix(statement, "INSERT") {
		values := strings.Split(statement, ",")
		if len(values) != len(headers) {
			match = false
			return
		}

		for _, value := range values {
			if value == "" {
				match = false
				return
			}
		}
	}

	// 检查UPDATE语句中列是否有效
	if strings.HasPrefix(statement, "UPDATE") {
		values := strings.Split(statement, ",")

		if len(values) < 3 {
			match = false
		}

		if _, ok := header2Index[values[len(values)-2]]; !ok {
			match = false
		}
	}

	return match
}
