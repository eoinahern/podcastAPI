package util

import "github.com/dchest/uniuri"

//GenerateRandomToken : get a  random token string
func GenerateRandomToken() string {
	return uniuri.NewLen(10)
}
