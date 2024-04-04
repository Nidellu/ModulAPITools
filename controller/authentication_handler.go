package controller

import (
	"encoding/json"
	"fmt"
	m "modulgo/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Kohansssgantengg123")
var tokenName = "token"

func generateToken(w http.ResponseWriter, id int, name string, userType int) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	claims := &m.Claims{
		ID:             id,
		Name:           name,
		UserType:       userType,
		StandardClaims: jwt.StandardClaims{ExpiresAt: tokenExpiryTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidtoken := validateUserToken(r, accessType)
		if !isValidtoken {
			sendAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType int) bool {
	isAccessTokenValid, id, email, userType := validateTokenFromCookies(r)
	fmt.Print(id, email, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, int, string, int) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		accessToken := cookie.Value
		accessClaims := &m.Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}

	return false, -1, "", -1
}

func sendAuthorizedResponse(w http.ResponseWriter) {
	var response m.UsersResponse
	response.Status = 401
	response.Message = "Error ga connect"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
