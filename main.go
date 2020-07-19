package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/controllers"
	_ "github.com/GhvstCode/Blog-Api/models"
	l "github.com/GhvstCode/Blog-Api/utils/logger"
)

func main(){
	handleRequest()
}

func handleRequest(){
	r := mux.NewRouter().StrictSlash(true)
	u := r.PathPrefix("/api/user").Subrouter()
	//b := r.PathPrefix("/blog").Subrouter()


	u.HandleFunc("/new", controllers.NewUser).Methods(http.MethodPost)
	u.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/resetPassword", controllers.ResetPassword).Methods(http.MethodPost)
	r.HandleFunc("/recoverPassword/{id}/{t}", controllers.RecoverPassword)
	//
	//b.HandleFunc("/", returnAllArticles)
	//b.HandleFunc("/new", createNewArticle).Methods("POST")
	//b.HandleFunc("/{id}", updateArticle).Methods("PUT")
	//b.HandleFunc("/{id}", deleteArticle).Methods("DELETE")
	//b.HandleFunc("/{id}", returnSingleArticle)
	//b.HandleFunc("/{id}/subscribe", returnSingleArticle)

	//log.Fatal(http.ListenAndServe(":8080", r))
	//l := log.New(os.Stdout, " product-api", log.LstdFlags)
	s := &http.Server{
		Addr: ":8080",
		Handler: r,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func () {
		l.InfoLogger.Println("Server is up on port",s.Addr)
		err := 	s.ListenAndServe()
		if err != nil {
			l.ErrorLogger.Fatal("Error starting server on port",s.Addr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.InfoLogger.Println("Received terminate, graceful shutdown! Signal: ",sig)

	tc, _:= context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)

}
///home/tobax/mongodb/bin/mongod --dbpath=/home/tobax/mongodb-data