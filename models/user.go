package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

//type Token struct {
//	UserId uint
//	jwt.StandardClaims
//}
//a struct to represent a User
type UserModel struct {
	Name 		string 	`json:"name"`
	Email    	string 	`json:"email"`
	Password 	string 	`json:"password"`
	Token    	string 	`json:"token";sql:"-"`
}
