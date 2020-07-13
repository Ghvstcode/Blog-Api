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
func Message(result bool, message string) map[string]interface{} {
	return map[string]interface{} {"result" : result, "message" : message}
}
//Testing Testing
//The hell yu waiting for??
func Response(result bool, message string, statuscode int) *Data {
	 return &Data{
	 	statusCode : statuscode,
	 	message : message,
	 	result : result,
	 }
}

