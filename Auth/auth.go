package Auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/GhvstCode/Blog-Api/models"
	"github.com/GhvstCode/Blog-Api/utils"
)

var Jwt = func(next http.Handler) http.Handler {
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		openRoutes := []string{"/api/user/new", "/api/user/login", "api/resetPassword", "api/recoverPassword", "/logs"}
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

		tokenPart := tArray[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("os.Getenv"), nil
		})

		if err != nil || !token.Valid{ //Malformed token, returns with http code 403 as usual
			utils.Response(false, "Malformed authentication token", http.StatusBadRequest).Send(w)
			return
		}
		i := strings.Split(tk.UserId, "_")
		UserID := i[0]
		//fmt.Print("This is the Users ID: ", UserID)
		ctx := context.WithValue(r.Context(), "user", UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
