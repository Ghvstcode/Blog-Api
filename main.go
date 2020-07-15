package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/controllers"
	_ "github.com/GhvstCode/Blog-Api/models"
)

func main(){
	fmt.Print("Hello")
	handleRequest()
}

func handleRequest(){
	r := mux.NewRouter().StrictSlash(true)
	u := r.PathPrefix("/api/user").Subrouter()
	//b := r.PathPrefix("/blog").Subrouter()


	u.HandleFunc("/new", controllers.NewUser).Methods(http.MethodPost)
	u.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	u.HandleFunc("/resetPassword", controllers.ResetPassword).Methods(http.MethodPost)
	//u.HandleFunc("/recoverPassword/{id}/{t}", returnAllArticles)
	//
	//b.HandleFunc("/", returnAllArticles)
	//b.HandleFunc("/new", createNewArticle).Methods("POST")
	//b.HandleFunc("/{id}", updateArticle).Methods("PUT")
	//b.HandleFunc("/{id}", deleteArticle).Methods("DELETE")
	//b.HandleFunc("/{id}", returnSingleArticle)
	//b.HandleFunc("/{id}subscribe", returnSingleArticle)

	//log.Fatal(http.ListenAndServe(":8080", r))
	l := log.New(os.Stdout, " product-api", log.LstdFlags)
	s := &http.Server{
		Addr: ":8080",
		Handler: r,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func () {
		l.Println("Server is up on port",s.Addr)
		err := 	s.ListenAndServe()
		if err != nil {
			l.Println("Error starting server on port",s.Addr)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Print("Received terminate, graceful shutdown! Signal: ",sig)

	tc, _:= context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)

}
///home/tobax/mongodb/bin/mongod --dbpath=/home/tobax/mongodb-data