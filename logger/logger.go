package logger

import (
	"encoding/json"
	"fmt"
)

// Payload is a structure for log
type Payload struct {
	Application  string `json:"app"`
	Method       string `json:"method"`
	Scheme       string `json:"scheme"`
	Origin       string `json:"origin"`
	Target       string `json:"target"`
	ResponseTime string `json:"response_time"`
}

// Do execute log event
func Do(p Payload) error {
	content, err := json.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	return nil
}
