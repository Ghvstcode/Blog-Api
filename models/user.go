package models

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/GhvstCode/Blog-Api/utils"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to represent a User
type UserModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name 		string 	`bson:"name" json:"name"`
	Email    	string 	`bson:"email" json:"email"`
	Password 	string 	`bson:"password" json:"password"`
	Token    	string 	`bson:"token" json:"token"`
	Subscriptions []Sub `bson:"sub, omitempty" json:"sub, omitempty"`
}

type Sub struct {
	SubID string `bson:"subID" json:"subID"`
}

func (u *UserModel) validate() map[string]interface{} {
	if !strings.Contains(u.Email, "@") {
		return utils.Message(false, "Provide a valid email address")
	}
	if passwordValid, _ :=regexp.MatchString("^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$", u.Password); !passwordValid{
		return utils.Message(false, "Password must be longer than 8 chars, and contain at least one digit!")
	}

	return utils.Message(true, "Validated")
}
func (u *UserModel) Create() map[string]interface{}{
	resp := u.validate()
	ok := resp["result"].(bool)
	if !ok {
		return  resp
	}
}