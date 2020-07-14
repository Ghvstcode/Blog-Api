package models

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/GhvstCode/Blog-Api/utils"
)

type Token struct {
	UserId string
	jwt.StandardClaims
}

//a struct to represent a User
type UserModel struct {
	_Id          primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	//_Id          string `json:"id"`
	ID            string             `json:"ID, omitempty"`
	Name          string             `bson:"name" json:"name"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"password"`
	Subscriptions []Sub              `bson:"sub, omitempty" json:"sub, omitempty"`
}

type Sub struct {
	SubID string `bson:"subID" json:"subID"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func genAuthToken(u *UserModel)(string, error){
	t := &Token{
		UserId: u.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), t)
	tokenString, err := token.SignedString([]byte("os.Getenv"))//Change this to load the jwt secret from env file.
	//tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *UserModel) validate() *utils.Data {
	if !strings.Contains(u.Email, "@") {
		return utils.Response(false, "Provide a valid email address", http.StatusBadRequest)
	}
	//if _, err :=regexp.MatchString(`(^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$)`, u.Password); err != nil{
	//	fmt.Print(err)
	//	return utils.Response(false, "Password must be longer than 8 chars, and contain at least one digit!" , http.StatusBadRequest)
	//}
	if len(u.Password) < 6 {
		return utils.Response(false, "Password is required", http.StatusBadRequest)
	}

	if u.Name == ""{
		return utils.Response(false, "Name is required", http.StatusBadRequest)
	}
	if len(u.Name) < 3 {
		return utils.Response(false, "Name is required", http.StatusBadRequest)
	}
	//Check to see if password is secure
	if strings.Contains(u.Password, "abcdefg") {
		return utils.Response(false, "Please provide a valid password", http.StatusBadRequest)
	}

	ErrorChan := make(chan error, 1)
	defer close(ErrorChan)
	go func(){
		ErrorChan <- User.FindOne(context.TODO(), bson.M{"email": u.Email}).Decode(u)
	}()
	Error := <- ErrorChan
	if Error == nil {
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
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}
	u.Password = string(hashedPassword)

	//res, err := User.InsertOne(context.TODO(), &u)
	//res, err := User.InsertOne(context.TODO(), &UserModel {
	//	//_Id:           primitive.NewObjectID(),
	//	Name:          u.Name,
	//	Email:         u.Email,
	//	Subscriptions: u.Subscriptions,
	//	Password:      u.Password,
	//	//ID: u._Id,
	//})
	if err != nil {
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			u.ID = oid.Hex()
	}

	t, e := genAuthToken(u)
	if e != nil {
		return utils.Response(false, "Failed to create account, connection error.", http.StatusBadGateway)
	}

	u.Password = ""
	response := utils.Response(true, "created", http.StatusCreated)
	response.Token = t
	response.Data = [1]*UserModel{u}
	return response
}
func Login (email string,  password string) *utils.Data {
	user := &UserModel{}

	ErrorChan := make(chan error, 1)
	defer close(ErrorChan)
	go func(){
		ErrorChan <- User.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)
	}()
	Error := <- ErrorChan
	fmt.Print(user)
	if Error != nil {
		fmt.Println(Error)
		return utils.Response(false, "Unable to  Login!" , http.StatusBadRequest)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utils.Response(false, "Invalid login credentials. Please try again", http.StatusUnauthorized)
	}
	user.Password = ""
	//user.ID = user._Id.Hex()
	t, _ := genAuthToken(user)//We would eventually check for the error & Log it later bla bla bla
	response := utils.Response(true, "created", http.StatusCreated)
	response.Token = t
	response.Data = [1]*UserModel{user}
	return response
}