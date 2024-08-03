package router

import (
	"github.com/gorilla/mux"
	"github.com/maulikk04/golang-database/controller"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/create", controller.CreateHandler).Methods("POST")
	r.HandleFunc("/read/{id}", controller.ReadHandler).Methods("GET")
	r.HandleFunc("/readall", controller.ReadAllHandler).Methods("GET")
	r.HandleFunc("/delete/{id}", controller.DeleteHandler).Methods("DELETE")

	return r
}
