package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/GhvstCode/Blog-Api/utils/logger"
)

type Data struct {
	StatusCode int
	Message string
	Result bool
	Data interface{}
	Token string `json:"Token, omitempty"`
}

func (data Data) Send(w http.ResponseWriter) interface{} {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	return json.NewEncoder(w).Encode(data)
}

//Message is exported
func Message(result bool, message string) map[string]interface{} {
	return map[string]interface{} {"result" : result, "message" : message}
}

func Response(result bool, message string, statusCode int) *Data {
	 return &Data{
	 	StatusCode : statusCode,
	 	Message : message,
	 	Result : result,
	 }
}

func Email(email string, Name string, Token string, Host string, Id string){
	fmt.Printf(Host+"/recoverPassword/"+ Id +"/"+Token )
	from := mail.NewEmail("BlogAPI", "BlogAPI@exaample.com")
	subject := "Password Reset"
	to := mail.NewEmail(Name, email)
	content := mail.NewContent("text/plain", "Click on the link below to reset the password for your Bloggy account\n " + Host+"/recoverPassword/"+ Id +"/"+Token + "\nThis link expires in 15 minutes. Ignore this mail if you had nothing to do with this.")
	m := mail.NewV3MailInit(from, subject, to, content)
	apiKey,ok := os.LookupEnv("SENDGRID_API_KEY")
	if ok == false{
		apiKey = os.Getenv("SENDGRID_API_KEY")
	}
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		logger.WarningLogger.Println("An Error occurred with the Email service")
		Response(false, "Unable to Send Password Rest", http.StatusInternalServerError)
		//log.Println(err)
	}
}

//Recover password =>
//curl --location --request POST 'http://localhost:8080/recoverPassword/5f1047d8ca88e9ff804e3376/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOm51bGwsImV4cCI6MTU5NDk5MDk0NH0.cjefmUpDadaM5pQKOBlSsNYOzg6dGLzrlErIXGt8gEk
//' \
//--header 'Content-Type: application/json' \
//--data-raw '{
// "name": "goldonboy",
// "email": "n9@example.com",
// "password": "ilovejesus9",
// "confirmpassword": "ilovejesus9"
//}'

//Reset Password =>
//curl --location --request POST 'http://localhost:8080/resetPassword' \
//--header 'Content-Type: application/json' \
//--data-raw '{
//  "name": "goldonboy",
//  "email": "nl9@example.com",
//  "password": "ilovejesus4"
//}'

//Login => curl --location --request POST 'http://localhost:8080/api/user/login' \
//--header 'Content-Type: application/json' \
//--data-raw '{
//    "name": "goldonboy",
//    "email": "nl9@example.com",
//    "password": "ilovejesus4"
//}'