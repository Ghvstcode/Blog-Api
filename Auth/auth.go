package Auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/GhvstCode/Blog-Api/models"
	"github.com/GhvstCode/Blog-Api/utils"
)

var Jwt = func(next http.Handler) http.Handler {
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		openRoutes := []string{"/api/user/new", "/api/user/login", "api/resetPassword", "api/recoverPassword"}
		requestPath := r.URL.Path

		for _, value := range openRoutes {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			utils.Response(false, "Missing auth token", http.StatusUnauthorized).Send(w)
			return
		}

		tArray := strings.Split(tokenHeader, " ")
		if len(tArray) != 2 {
			utils.Response(false, "Invalid/Malformed Auth token", http.StatusUnauthorized).Send(w)
			return
		}

		tokenPart := tArray[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("os.Getenv"), nil
		})
	})
}
