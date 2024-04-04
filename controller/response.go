package controller

import (
	"encoding/json"
	m "modulgo/model"
	"net/http"
)

func sendUserSuccessResponse(w http.ResponseWriter, message string, users []m.User) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

func sendUserErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 400
	response.Message = message
	json.NewEncoder(w).Encode(response)
}
