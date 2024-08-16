package models

import "encoding/json"

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	json, _ := json.Marshal(e)
	return string(json[:])
}
