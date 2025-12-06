package utils

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
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