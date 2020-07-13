package models

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

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

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *UserModel) validate() *utils.Data {
	if !strings.Contains(u.Email, "@") {
		return utils.Response(false, "Provide a valid email address", http.StatusBadRequest)
	}
	if passwordValid, _ :=regexp.MatchString("^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$", u.Password); !passwordValid{
		return utils.Response(false, "Password must be longer than 8 chars, and contain at least one digit!" , http.StatusBadRequest)
	}
	_, err := User.Find(context.TODO(), bson.D{{"name", "Bob"}})
	if err == nil {
		return utils.Response(false, "Invalid Email!" , http.StatusBadRequest)
	}

	return utils.Response(true, "Validated", http.StatusAccepted)
}
func (u *UserModel) Create() *utils.Data{
	resp := u.validate()
	ok := resp.Result; if !ok {
		return resp
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return utils.Response(false, "An error occurred! Unable to save user", http.StatusInternalServerError)
	}
	u.Password = string(hashedPassword)
}