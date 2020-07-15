package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Data struct {
	StatusCode int
	Message string
	Result bool
	Data interface{}
	Token string
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
	fmt.Printf("Click on the link below to reset the password for your Bloggy account\n " + Host+"/recoverPassword/"+ Id +"/"+Token + "\nThis link expires in 15 minutes. Ignore this mail if you had nothing to do with this.")
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
		Response(false, "Unable to Send Password Rest", http.StatusInternalServerError)
		//log.Println(err)
	}
}
