package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GhvstCode/Blog-Api/models"
	"github.com/GhvstCode/Blog-Api/utils"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value("user").(string)
	//Content goes  here!
	post := &models.BlogModel{}

	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		utils.Response(false, "An error occurred, Unable to create post", http.StatusBadRequest)
	}

	resp := post.Create(ownerID)
	resp.Send(w)
}

func UpdatePost (w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	//Content goes  here!
	post := &models.UpdateBlogModel{}

	err := json.NewDecoder(r.Body).Decode(post)
	fmt.Print(post)
	if err != nil {
		utils.Response(false, "An error occurred, Unable to create post", http.StatusBadRequest)
	}

	resp := post.UpdatePost(id)
	resp.Send(w)
}

func DeletePost(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user").(string)
	vars := mux.Vars(r)
	id := vars["id"]
	resp :=models.DeletePost(id)
	resp.Send(w)
}