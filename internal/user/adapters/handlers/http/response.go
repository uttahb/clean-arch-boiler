package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message" omitempty:"true"`
	Data    interface{} `json:"data" omitempty:"true"`
}

func SuccessResponse(data interface{}, message string) *Response {
	response := &Response{
		Ok:      true,
		Message: message,
		Data:    data,
	}
	return response
}
func ErrorResponse(message string) *Response {
	response := &Response{
		Ok:      false,
		Message: message,
	}
	return response
}
func (resp *Response) Send(w http.ResponseWriter, r *http.Request, status int) {
	render.Status(r, status)
	render.JSON(w, r, resp)
}
