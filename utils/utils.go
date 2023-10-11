package utils

import (
	"encoding/json"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	//marshData, err := json.Marshal(data)
	//if err != nil {
	//	logrus.Errorf("EncodeJSONBody : Error marshing response data interface %v", err)
	//}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("EncodeJSONBody : Error encoding response %v", err)
	}
}

func EncodeJSON200Body(resp http.ResponseWriter, data interface{}) {
	var newJSON = jsoniter.ConfigCompatibleWithStandardLibrary
	err := newJSON.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("EncodeJSON200Body : Error encoding response %v", err)
	}
}
