package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/eoinahern/podcastAPI/models"

	"github.com/eoinahern/podcastAPI/util"
)

/* stole this middleware implementation from https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
** not sure how intuitive it is to look at but works so ill keep it for now.
 */

//Adapter : function type adapter
type Adapter func(http.Handler) http.Handler

//Adapt : wrap middleware adapters around finally executed route code
func Adapt(finalHandler http.Handler, adapters ...Adapter) http.Handler {

	for _, item := range adapters {
		finalHandler = item(finalHandler)
	}

	return finalHandler
}

//AuthMiddlewareInit : initial middleware executed to check jwt token is valid
func AuthMiddlewareInit(jwtTokenUtil *util.JwtTokenUtil) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			fmt.Println("authorization middleware!!")
			token := getTokenFromHeader(req)
			code, message := jwtTokenUtil.CheckTokenCredentials(token)

			if code != -1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				msg, _ := json.Marshal(models.Message{Message: message})
				w.Write(msg)
				fmt.Println("auth failed!!")
				return
			}

			h.ServeHTTP(w, req)

		})
	}
}

//PagingParamsValidate check paging values passed in request
func PagingParamsValidate(max uint16) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			queryParams := req.URL.Query()

			i, err := strconv.ParseUint(queryParams.Get("limit"), 10, 16)

			if err != nil {
				incorrectParamsError(w)
				return
			}

			limit := uint16(i)
			i, err = strconv.ParseUint(queryParams.Get("offset"), 10, 16)

			if err != nil {
				incorrectParamsError(w)
				return
			}
			offset := uint16(i)

			if limit > max || offset%10 != 0 || limit%10 != 0 {
				incorrectParamsError(w)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}

func incorrectParamsError(resp http.ResponseWriter) {
	http.Error(resp, "incorrect params", http.StatusBadRequest)
}

func getTokenFromHeader(req *http.Request) string {

	token := req.Header.Get("Authorization")
	tokenSlice := strings.Split(token, " ")

	if len(tokenSlice) != 2 {
		return ""
	}

	return tokenSlice[1]
}
