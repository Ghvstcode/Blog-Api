package controllers

import (
	"encoding/json"
	"net/http"

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