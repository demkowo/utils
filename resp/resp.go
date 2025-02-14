package resp

import (
	"net/http"
)

type Err struct {
	Error  string        `json:"error,omitempty"`
	Code   int           `json:"code"`
	Status string        `json:"status"`
	Causes []interface{} `json:"causes"`
}

type Ok struct {
	Message string        `json:"message,omitempty"`
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Data    []interface{} `json:"data,omitempty"`
}

func Error(code int, err string, causes []interface{}) *Err {
	return &Err{
		Error:  err,
		Code:   code,
		Status: http.StatusText(code),
		Causes: causes,
	}
}

func New(code int, message string, data []interface{}) *Ok {
	return &Ok{
		Message: message,
		Code:    code,
		Status:  http.StatusText(code),
		Data:    data,
	}
}

func (e *Err) JSON() (int, map[string]any) {
	code := e.Code
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return code, map[string]any{
		"error":  e.Error,
		"status": http.StatusText(code),
		"code":   code,
		"causes": e.Causes,
	}
}

func (o *Ok) JSON() (int, map[string]any) {
	code := o.Code
	if code == 0 {
		code = http.StatusOK
	}
	return code, map[string]any{
		"message": o.Message,
		"data":    o.Data,
		"status":  http.StatusText(code),
		"code":    code,
	}
}
