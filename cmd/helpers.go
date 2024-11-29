package main

import (
	"encoding/json"
	"net/http"
)

type httpResp struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func sendResponse(rw http.ResponseWriter, status int, data interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	out, err := json.Marshal(httpResp{Status: status, Data: data, Message: message})
	if err != nil {
		sendErrorResponse(rw, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	rw.Write(out)
}

func sendErrorResponse(rw http.ResponseWriter, status int, data interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	out, _ := json.Marshal(httpResp{Status: status,
		Message: message,
		Data:    data})

	rw.Write(out)
}
