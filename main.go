package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/api/Auth"

	"github.com/GhvstCode/Blog-Api/api/controllers"
	_ "github.com/GhvstCode/Blog-Api/api/models"
	l "github.com/GhvstCode/Blog-Api/api/utils/logger"
)

func main() {
	handleRequest()
}

func handleRequest() {
	r := mux.NewRouter().StrictSlash(true)
	u := r.PathPrefix("/api/user").Subrouter()
	b := r.PathPrefix("/api/blog").Subrouter()

	r.Use(Auth.Jwt)
	r.HandleFunc("/logs", controllers.ViewLog)
	u.HandleFunc("/new", controllers.NewUser).Methods(http.MethodPost)
	u.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/resetPassword", controllers.ResetPassword).Methods(http.MethodPost)
	r.HandleFunc("/api/recoverPassword/{id}/{t}", controllers.RecoverPassword)
	u.HandleFunc("/posts", controllers.GetPosts).Methods(http.MethodGet)
	//u.HandleFunc("/feed", controllers.GetPosts).Methods(http.MethodGet)This includes a list of all articles/A users subscriptions.

	b.HandleFunc("/new", controllers.NewPost).Methods("POST")
	b.HandleFunc("/{id}", controllers.UpdatePost).Methods("PUT")
	b.HandleFunc("/{id}", controllers.DeletePost).Methods("DELETE")
	b.HandleFunc("/{id}", controllers.GetOnePost).Methods(http.MethodGet)
	//b.HandleFunc("/{id}/subscribe", returnSingleArticle)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.InfoLogger.Println("Server is up on port", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.ErrorLogger.Fatal("Error starting server on port", s.Addr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.InfoLogger.Println("Received terminate, graceful shutdown! Signal: ", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)

}
