package backend

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		ErrorMsg string `json:"error"`
	}

	respMsg := errorResponse{ErrorMsg: msg}
	data, err := json.Marshal(respMsg)
	if err != nil {
		code = 500
		respMsg = errorResponse{ErrorMsg: "Something went wrong"}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
