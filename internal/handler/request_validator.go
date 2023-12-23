package handler

import (
	"fmt"
	"regexp"
	"strings"

	"zh.com/ms_coding2/internal/repository"
)

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
		if len(values) != len(repository.Headers) {
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

		if _, ok := repository.Header2Index[values[len(values)-2]]; !ok {
			match = false
		}
	}

	return match
}

func isValidQuery(query string) bool {
	// 定义正则
	pattern := ""
	for i, n := range repository.Headers {
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
