package utils

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	statusCode int
	message string
	result bool
}

func (data Data) send(w http.ResponseWriter) interface{} {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.statusCode)
	return json.NewEncoder(w).Encode(data)
}

//Message is exported
//func NewMessage(result bool, message string) map[string]interface{} {
//	return map[string]interface{} {"result" : result, "message" : message}
//}

func Response(result bool, message string, statuscode int) *Data {
	 return &Data{
	 	statusCode : statuscode,
	 	message : message,
	 	result : result,
	 }
}

