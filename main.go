package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/GhvstCode/Blog-Api/models"
)

func main(){
	fmt.Print("Hello")
	//handleRequest()
}

func handleRequest(){
	r := mux.NewRouter().StrictSlash(true)
	//u := r.PathPrefix("/user").Subrouter()
	//b := r.PathPrefix("/blog").Subrouter()


	//u.HandleFunc("/new", returnAllArticles)
	//u.HandleFunc("/login", returnAllArticles)
	//u.HandleFunc("/resetPassword", returnAllArticles)
	//u.HandleFunc("/recoverPassword", returnAllArticles)
	//
	//b.HandleFunc("/", returnAllArticles)
	//b.HandleFunc("/new", createNewArticle).Methods("POST")
	//b.HandleFunc("/{id}", updateArticle).Methods("PUT")
	//b.HandleFunc("/{id}", deleteArticle).Methods("DELETE")
	//b.HandleFunc("/{id}", returnSingleArticle)
	//b.HandleFunc("/{id}subscribe", returnSingleArticle)

	log.Fatal(http.ListenAndServe(":8080", r))
}
///home/tobax/mongodb/bin/mongod --dbpath=/home/tobax/mongodb-data