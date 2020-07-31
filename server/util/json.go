package util

import (
	"bufio"
	"encoding/json"
	"net/http"
)

// Respond forms a JSON response based on the given data
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ExtractJSON extracts JSON from a Request body if possible
func ExtractJSON(r *http.Request) (map[string]interface{}, error) {
	msg := make(map[string]interface{})

	if err := json.NewDecoder(bufio.NewReader(r.Body)).Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// Message is used to send a standard (normally failure) response message
func Message(w http.ResponseWriter, status bool, msg string) {
	Respond(w, map[string]interface{}{"status": status, "message": msg})
}
