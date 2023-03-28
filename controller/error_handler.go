package controller

import (
	"PBP-API-Tools-1121004-1121008-1121018-1121032/model"
	"encoding/json"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response model.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
