package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/models"
	"github.com/GhvstCode/Blog-Api/utils"
)

func NewUser (w http.ResponseWriter, r *http.Request) {
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res:= user.Create()
	res.Send(w)
}

func Login (w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Host)
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res := models.Login(user.Email, user.Password)
	res.Send(w)
}

func ResetPassword (w http.ResponseWriter, r *http.Request) {
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res := models.ResetPassword(user.Email, r.Host)
	res.Send(w)
}

func RecoverPassword(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	t := vars["t"]

	user := &models.RePassword{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
}
