package util

import (
	"log"
	"net/http"
	"time"

	"github.com/eoinahern/podcastAPI/repository"

	jwt "github.com/dgrijalva/jwt-go"
)

//JwtTokenUtilInt interface
type JwtTokenUtilInt interface {
	CreateToken(username string) string
	CheckTokenCredentials(token string) (int, string)
}

//JwtTokenUtil : helper methods for dealing with jwt token
type JwtTokenUtil struct {
	SigningKey string
	DB         repository.UserDBInt
}

//CreateToken : create a jwt token
func (j *JwtTokenUtil) CreateToken(username string) string {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["exp"] = time.Now().Add(time.Hour + 1).Unix()

	signedToken, err := token.SignedString([]byte(j.SigningKey))

	if err != nil {
		log.Fatal(err)
	}

	return signedToken
}

//CheckTokenCredentials : check token contains email in our records. plus check its not expired
func (j *JwtTokenUtil) CheckTokenCredentials(tokenStr string) (int, string) {

	token, err := jwt.Parse(tokenStr, func(passedToken *jwt.Token) (interface{}, error) {
		return []byte(j.SigningKey), nil
	})

	if err != nil {
		return http.StatusUnauthorized, "error validating token"
	}

	claims := token.Claims.(jwt.MapClaims)
	time := int64(claims["exp"].(float64))
	name := claims["name"].(string)

	if !verifyTokenTime(time) || !j.verifyTokenUser(name) {
		return http.StatusUnauthorized, "error validating token"
	}

	return -1, ""
}

func verifyTokenTime(chimey int64) bool {
	if chimey > time.Now().Unix() {
		log.Println("token isnt expired")
		return true
	}

	log.Println("token is expired")
	return false
}

func (j *JwtTokenUtil) verifyTokenUser(tokenName string) bool {

	if !j.DB.CheckExist(tokenName) {
		return false
	}

	log.Println("name comparison succeed")
	return true
}
