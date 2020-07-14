package controllers

import (
	"encoding/json"
	"net/http"

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
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res := models.Login(user.Email, user.Password)
	res.Send(w)
}
