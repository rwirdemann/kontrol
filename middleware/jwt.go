package middleware

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// This is just an example JWT validation that does the validation but serves the request even
// if validation fails
func JWTMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header != "" && len(strings.Split(header, " ")) == 2 {
			jwt := strings.Split(header, " ")[1]
			validate(jwt)
		}
		next.ServeHTTP(w, r)
	})
}

// Token taken from Keycloak RSA key
const rsaPublicKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAl+hgn1UL8XjGZA6PHJfA16+adPHcgFBReXqiw0ykWivX46Ass0Psw60Q6nxhgOU4fa0QNVuxeTrvxpgiPu2MhrSSJxsNW2BiYETfAdjCpbwqmR1JhpyL1UR1MOdtkIe+Ucy6tYmpL4lt9gkREgDpv8pfQXAk5tYlaGnhPZM/53kRV3N1cFYlYC65ykY+JDkJdT74gFKbekOtYQiJPfmeuBOtBNZ1FiSf7T9k1bhtks6Q9ZbDjr8y9ax5OHCZEJLhTzQTIN4YdpV5nFR3eBZ5/kOS/E60JUjTWKgMYuSeoRi5drcYxQXPlM/gmCgY9igKfNa4gEeGx1cq1LSxV01tQQIDAQAB
-----END RSA PRIVATE KEY-----`

func validate(tokenString string) bool {
	key, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaPublicKeyPEM))
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
		return true
	}
}
