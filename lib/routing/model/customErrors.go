package routingModel

import "time"

type ErrorResponse struct {
	Message    string
	ServerTime time.Time
}

func (er *ErrorResponse) Error() string {
	return er.Message
}

func NewErrorMessage(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		ServerTime: time.Now(),
	}
}

func NewError(err error) *ErrorResponse {
	return &ErrorResponse{
		Message:    err.Error(),
		ServerTime: time.Now(),
	}
}
