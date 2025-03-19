package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"Success"`
	Error   string `json:"Error,omitempty"`
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, errorCode int, error string) {
	render.JSON(w, r, ErrorResponse{
		Success: false,
		Error:   error,
	})
	w.WriteHeader(errorCode)
}
