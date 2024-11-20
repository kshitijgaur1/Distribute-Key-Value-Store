package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// In-memory key-value store
var store = make(map[string]string)

// SetKey handles setting a key-value pair
func SetKey(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	key, value := data["key"], data["value"]
	store[key] = value
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetKey handles retrieving the value for a given key
func GetKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value, exists := store[key]
	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}
