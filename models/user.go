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
	l "github.com/GhvstCode/Blog-Api/utils/logger"
)

type Token struct {
	UserId string
	jwt.StandardClaims
}

//a struct to represent a User
type UserModel struct {
	ID         primitive.ObjectID `bson:"_id, omitempty" json:"id, omitempty"`
	Name          string `bson:"name" json:"name"`
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password" json:"password"`
	Subscriptions []Sub  `bson:"sub, omitempty" json:"sub, omitempty"`
}

type ReUserModel struct {
	ID           string `json:"id, omitempty"`
	Name          string `bson:"name" json:"name"`
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password" json:"password"`
	Subscriptions []Sub  `bson:"sub, omitempty" json:"sub, omitempty"`
}

type RecPassword struct {
	Password           string `json:"password"`
	ConfirmPassword    string `json:"confirmPassword"`
}

type ResPassword struct {
	Email         string `json:"email"`
}

type Sub struct {
	SubID string `bson:"subID" json:"subID"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func genAuthToken(id string)(string, error){
	t := &Token{
		UserId: id + "_" + time.StampMicro,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), t)
	tokenString, err := token.SignedString([]byte("os.Getenv"))//Change this to load the jwt secret from env file.
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
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}
	u.Password = string(hashedPassword)


	res, err := User.InsertOne(context.TODO(), &UserModel {
		ID:           primitive.NewObjectID(),
		Name:          u.Name,
		Email:         u.Email,
		Subscriptions: u.Subscriptions,
		Password:      u.Password,
	})

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}

	var UID string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			UID = oid.Hex()
	}

	t, e := genAuthToken(UID)
	if e != nil {
		l.ErrorLogger.Println(e)
		return utils.Response(false, "Failed to create account, connection error.", http.StatusBadGateway)
	}

	u.Password = ""
	v := &ReUserModel{
		ID:            UID,
		Name:          u.Name,
		Email:         u.Email,
		Password:      "",
		Subscriptions: u.Subscriptions,
	}
	response := utils.Response(true, "created", http.StatusCreated)
	response.Token = t
	response.Data = [1]*ReUserModel{v}
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
	if Error != nil {
		l.ErrorLogger.Println(Error)
		return utils.Response(false, "Unable to  Login!" , http.StatusBadRequest)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Invalid login credentials. Please try again", http.StatusUnauthorized)
	}
	user.Password = ""
	//user.ID = user._Id.Hex()
	t, _ := genAuthToken(user.ID.Hex())//We would eventually check for the error & Log it later bla bla bla
	fmt.Print("LoginTok", user.ID.Hex())
	response := utils.Response(true, "Logged In", http.StatusOK)
	response.Token = t
	response.Data = [1]*UserModel{user}
	return response
}

func ResetPassword(email string, Host string) *utils.Data {
	u := &UserModel{}
	//1.Take in the users Email!
	//2.Check the User exists in the DB and pass the users ID!
	ErrorChan := make(chan error, 1)
	defer close(ErrorChan)
	go func(){
		ErrorChan <- User.FindOne(context.TODO(), bson.M{"email": email}).Decode(u)
	}()
	Error := <- ErrorChan
	if Error != nil {
		return utils.Response(false, "An Error occurred" , http.StatusBadRequest)
	}
	//3. Generate the auth token using the UsersName & Email along with the created at field as its signature!
	t := &Token{
		UserId: u.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), t)
	tokenString, err := token.SignedString([]byte(u.ID.String()+"."+u.Email))
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred" , http.StatusInternalServerError)
	}

	utils.Email(u.Email, u.Name, tokenString, Host, u.ID.Hex())
	return utils.Response(true, "Email Sent" , http.StatusAccepted)
}

func RecoverPassword(P *RecPassword, id string, tkn string) *utils.Data{
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to recover password" , http.StatusInternalServerError)
	}

	u := &UserModel{}

	ErrorChan := make(chan error, 1)
	defer close(ErrorChan)
	go func(){
		ErrorChan <- User.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(u)
	}()
	Error := <- ErrorChan
	if Error != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unable to recover password" , http.StatusInternalServerError)
	}

	tk := &Token{}

	token, err := jwt.ParseWithClaims(tkn, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.ID.String()+"."+u.Email), nil
	})

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Malformed authentication token", http.StatusForbidden)
	}

	if !token.Valid {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Malformed authentication token", http.StatusForbidden)
	}

	if P.Password != P.ConfirmPassword{
		l.ErrorLogger.Println("Password mismatch")
		return utils.Response(false, "Passwords do not match", http.StatusBadRequest)
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to Change password", http.StatusInternalServerError)
	}

	filter := bson.D{{"email", u.Email}}
	update := bson.D{{"$set", bson.D{{"password", hashedPassword}}}}

	_, err = User.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to Change password", http.StatusInternalServerError)
	}

	u.Password = ""
	response := utils.Response(true, "created", http.StatusCreated)
	response.Data = [1]*UserModel{u}
	return response
}

func GetPosts(id string) *utils.Data {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Invalid ID" , http.StatusInternalServerError)
	}

	post := &BlogModel{}

	cursor, err := Blog.Find(context.TODO(), bson.M{"owner": oid})

	if cursor == nil {
		return utils.Response(false, "You do not have any BlogPosts" , http.StatusBadRequest)
	}

	if err = cursor.All(context.TODO(), &post); err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Unable to  Fetch post" , http.StatusBadRequest)
	}

	response := utils.Response(true, "Logged In", http.StatusOK)
	response.Data = []*BlogModel{post}
	return response
}
