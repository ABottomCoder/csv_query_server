package server

import (
	"fmt"
	"net/http"
	"regexp"
)

func RegisterQueryHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", queryHandler)
}

type QueryResponse struct {
	Result [][]string `json:"result"`
	Msg    string     `json:"msg"`
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	fmt.Printf("query: %s, %v\n", query, isValidQuery(query))

	// 检查查询是否有效
	if isValidQuery(query) {
		// 处理查询并返回结果
		result := executeQuery(query)
		response := QueryResponse{
			Result: result,
		}
		jsonResponse(w, http.StatusOK, response)
	} else {
		// 返回查询格式错误的响应
		response := QueryResponse{
			Msg: "Error description of the malformed query",
		}
		jsonResponse(w, http.StatusBadRequest, response)
	}
}

func isValidQuery(query string) bool {
	// 定义正则
	pattern := ""
	for i, n := range headers {
		if i > 0 {
			pattern += "|"
		}
		pattern += regexp.QuoteMeta(n)
	}

	// fmt.Printf("pattern: %s\n", pattern)
	// validQueryRegex := `^((C1|C2|C3)\s*(==|!=|\$=|&=)\s*.*\s*(and|or)?\s*)+$`
	// validQueryRegex := `^(([A-Za-z0-9]+)\s*(==|!=|\$=|&=)\s*.*\s*(and|or)?\s*)+$`
	// validQueryRegex := `^(([A-Za-z0-9]+)\s*(==|!=|\$=|&=)\s*(?=\S).*(and|or)?\s*)+$`
	// validQueryRegex := `^(([A-Za-z0-9]+)\s*(==|!=|\$=|&=)\s*\S+\s*(and|or)?\s*)+$`
	// validQueryRegex := `^((C1|C2|C3)\s*(==|!=|\$=|&=)\s*\S+\s*(and|or)?\s*)+$`
	validQueryRegex := fmt.Sprintf("^((%s)\\s*(==|!=|\\$=|&=)\\s*\\S+\\s*(and|or)?\\s*)+$", pattern)

	// 匹配查询语句
	match, _ := regexp.MatchString(validQueryRegex, query)

	return match
}

// 判断切片中是否包含指定元素
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
