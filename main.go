package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	s := route.PathPrefix("/api").Subrouter()

	s.HandleFunc("/createprofile", createProfile).Methods("POST")
	s.HandleFunc("getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/get/{id}", getUserProfile).Methods("GET")
	s.HandleFunc("/updateprofile", updateProfile).Methods("PUT")
	s.HandleFunc("/deleteprofile", deleteProfile).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9080", s))
}
