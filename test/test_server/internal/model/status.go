package model

import "time"

type Status struct {
	Message    string    `json:"message"`
	ServerTime time.Time `json:"server_time"`
}

func NewStatus(message string) *Status {
	return &Status{
		Message:    message,
		ServerTime: time.Now(),
	}
}
