package main

import (
	"fmt"
	"net/http"

	"github.com/maulikk04/golang-database/controller"
	"github.com/maulikk04/golang-database/model"
	"github.com/maulikk04/golang-database/router"
)

func main() {
	dir := "./"
	db, err := model.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	controller.Initialize(db)
	r := router.SetupRouter()

	fmt.Println("Listening on Server 4000")

	if err := http.ListenAndServe(":4000", r); err != nil {
		fmt.Println("Error", err)
	}
}
