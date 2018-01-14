package middleware

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// This is just an example JWT validation that does the validation but serves the request even
// if validation fails
func JWTMiddleware(keycloakRSAPub []byte, next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header != "" && len(strings.Split(header, " ")) == 2 {
			tokenString := strings.Split(header, " ")[1]
			validate(keycloakRSAPub, tokenString)
		}
		next.ServeHTTP(w, r)
	})
}

func validate(keycloakRSAPub []byte, tokenString string) bool {
	key, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(keycloakRSAPub))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["name"])
		return true
	} else {
		fmt.Println(err)
		return false
	}
}
