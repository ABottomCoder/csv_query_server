package handler

import (
	"fmt"
	"net/http"

	"zh.com/ms_coding2/internal/repository"
	"zh.com/ms_coding2/internal/server"
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
		err := repository.ExecuteModify(job)
		if err != nil {
			response.Msg = "Execute modify fail"
			fmt.Printf("Execute modify fail, err: %v\n", err)
			server.JsonResponse(w, http.StatusBadRequest, response)
		} else {
			response.Msg = "Modification successful"
			server.JsonResponse(w, http.StatusOK, response)
		}
	} else {
		response := ModifyResponse{
			Msg: "Error description of the malformed modify",
		}
		server.JsonResponse(w, http.StatusBadRequest, response)
	}
}
