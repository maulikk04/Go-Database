package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/maulikk04/golang-database/model"
)

var db *model.Driver

func Initialize(database *model.Driver) {
	db = database
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	data["id"] = id

	if err := db.Write("users", id, data); err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	data, err := db.Read("users", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func ReadAllHandler(w http.ResponseWriter, r *http.Request) {

	data, err := db.ReadAll("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	if err := db.Delete("users", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
