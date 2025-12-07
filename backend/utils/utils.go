package utils

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"errors"
)
//encrypt password
func HashPassword(pwd string) (string,error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd),12)
	return string(hash),err
}
// generate json web token
func GenerateJWT(username string)(string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"username" : username,
		"exp" : time.Now().Add(time.Hour*72).Unix(),
	})

	signedToken, err := token.SignedString([]byte("secret"))
	return "bearer "+signedToken, err
}
//check password
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {return true} else{return false}
}

func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "bearer "{
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signin method")
		}
		return []byte("secret"), nil
	})

	if err!=nil {
		return "" ,err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		username , ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username claim is not a string")
		}
		return username, nil
	}
	
	return "" ,err
}