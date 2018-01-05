package main

import (
	"net/http"
	"demo/controller"
	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/users", controller.G_usc.ListAllUsers).Methods("GET")
	r.HandleFunc("/users", controller.G_usc.CreateUser).Methods("POST")
	r.HandleFunc("/users/{user_id}/relationships", controller.G_usc.ListAllRelationshipOfUser).Methods("GET")
	r.HandleFunc("/users/{user_id}/relationships/{other_user_id}", controller.G_usc.UpdateRelationship).Methods("PUT")

	http.ListenAndServe(":8080", r)
}