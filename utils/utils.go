package utils

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	StatusCode int
	Message string
	Result bool
	Data interface{}
}

func (data Data) Send(w http.ResponseWriter) interface{} {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	return json.NewEncoder(w).Encode(data)
}

//Message is exported
func Message(result bool, message string) map[string]interface{} {
	return map[string]interface{} {"result" : result, "message" : message}
}

func Response(result bool, message string, statuscode int) *Data {
	 return &Data{
	 	StatusCode : statuscode,
	 	Message : message,
	 	Result : result,
	 }
}

