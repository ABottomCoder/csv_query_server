package handler

import (
	"fmt"
	"net/http"

	"zh.com/ms_coding2/internal/repository"
	"zh.com/ms_coding2/internal/server"
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
		result := repository.ExecuteQuery(query)
		response := QueryResponse{
			Result: result,
		}
		server.JsonResponse(w, http.StatusOK, response)
	} else {
		// 返回查询格式错误的响应
		response := QueryResponse{
			Msg: "Error description of the malformed query",
		}
		server.JsonResponse(w, http.StatusBadRequest, response)
	}
}
