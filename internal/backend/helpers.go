package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		ErrorMsg string `json:"error"`
	}

	respMsg := errorResponse{ErrorMsg: msg}
	data, err := json.Marshal(respMsg)
	if err != nil {
		code = 500
		respMsg = errorResponse{ErrorMsg: serverErrorString}
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

func replaceBadWords(opinion string) (string, error) {
	badWordsMap, err := readJsonFile(badWordsFilePath)
	if err != nil {
		return "", err
	}
	badWords := badWordsMap["bad_words"]

	opinionList := strings.Split(opinion, " ")
	newOpinionList := make([]string, 0, len(opinionList))
	for _, word := range opinionList {
		badWordFound := false
		for _, badWord := range badWords {
			if strings.ToLower(badWord) == strings.ToLower(word) {
				badWordFound = true
				break
			}
		}
		if badWordFound {
			newOpinionList = append(newOpinionList, "****")
		} else {
			newOpinionList = append(newOpinionList, word)
		}
	}
	return strings.Join(newOpinionList, " "), nil
}

func readJsonFile(filePath string) (map[string]([]string), error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string][]string
	err = json.Unmarshal(f, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
