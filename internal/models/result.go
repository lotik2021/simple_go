package models

import (
	"encoding/json"
	"fmt"
)

type UapiResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error UapiError   `json:"errors,omitempty"`
}

type RawResponse struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error UapiError       `json:"errors,omitempty"`
}

type UapiError []interface{}

func (uErr UapiError) Error() string {
	return fmt.Sprintf("%#v", uErr)
}
