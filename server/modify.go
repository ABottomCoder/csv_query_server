package server

import (
	"net/http"
	"regexp"
	"strings"
)

func RegisterModifyHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", modifyHandler)
}

// http://localhost:7259/?job=INSERT%20a1%2Ca2%2Ca3
// http://localhost:7259/?job=DELETE%20a1
// http://localhost:7259/?job=DELETE%20a1%2Ca2
// http://localhost:7259/?job=UPDATE%20ab%2Cac%2CC2%2Caa
func modifyHandler(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")

	if job == "" {
		http.Error(w, "Missing job parameter", http.StatusBadRequest)
		return
	}

	err := executeModify(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Modification successful"))
}

func isValidModify(statement string) bool {
	// 定义有效语句的正则表达式
	validStatementRegex := `^(INSERT|DELETE|UPDATE)\s+".*"(,\s*".*")*$`

	// 使用正则表达式匹配输入语句
	match, _ := regexp.MatchString(validStatementRegex, statement)

	// 检查INSERT语句中是否包含3个列的值
	if strings.HasPrefix(statement, "INSERT") {
		values := strings.Split(statement, ",")
		if len(values) != 3 {
			match = false
		}
	}

	return match
}
