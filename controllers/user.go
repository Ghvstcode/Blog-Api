package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/models"
	"github.com/GhvstCode/Blog-Api/utils"
	l "github.com/GhvstCode/Blog-Api/utils/logger"
)

func NewUser (w http.ResponseWriter, r *http.Request) {
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res:= user.Create()
	res.Send(w)
}

func Login (w http.ResponseWriter, r *http.Request) {
	//fmt.Print(r.Host)
	user := &models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res := models.Login(user.Email, user.Password)
	res.Send(w)
}

func ResetPassword (w http.ResponseWriter, r *http.Request) {
	user := &models.ResPassword{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
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

	user := &models.RecPassword{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}

	res := models.RecoverPassword(user, id, t)
	res.Send(w)
}

func GetPosts(w http.ResponseWriter, r *http.Request){
	ID := r.Context().Value("user").(string)
	resp := models.GetPosts(ID)
	resp.Send(w)
}

func ViewLog(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("logs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	b, err := ioutil.ReadAll(file)
	_, _ = w.Write(b)
}

//func feed(w http.ResponseWriter, r *http.Request){
//	//var posts []*BlogModel
//}