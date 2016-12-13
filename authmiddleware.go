package authmiddleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/authmiddleware/bearertoken"
	"github.com/byuoitav/authmiddleware/wso2jwt"
	"github.com/jessemillar/jsonresp"
)

// Authenticate is the middleware function
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Go down the "if chain" checking for various authentication methods
		checkLocal(next, writer, request)
		checkBearerToken(next, writer, request)
		checkWSO2(next, writer, request)
	})
}

func checkLocal(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	log.Printf("Local check starting")
	if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
		log.Printf("Local check finished")
		next.ServeHTTP(writer, request)
		return
	}
}

func checkBearerToken(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	log.Printf("Bearer tokenk check starting")
	token := request.Header.Get("Authorization") // Get the token if it exists

	if len(token) > 0 { // Proceed if we found a token
		parts := strings.Split(token, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonresp.New(writer, http.StatusBadRequest, "Bad Authorization header")
			return
		}

		valid, err := bearertoken.CheckToken([]byte(parts[1])) // Validate the existing token
		if err != nil {
			jsonresp.New(writer, http.StatusBadRequest, err.Error())
			return
		}

		if valid {
			log.Println("Bearer token authorized")
			next.ServeHTTP(writer, request)
			return
		}
	}
}

func checkWSO2(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	log.Printf("WSO2 check starting")
	token := request.Header.Get("X-jwt-assertion") // Get the token if it exists

	if len(token) > 0 { // Proceed if we found a token
		valid, err := wso2jwt.Validate(token) // Validate the existing token
		if err != nil {
			jsonresp.New(writer, http.StatusBadRequest, err.Error())
			return
		}

		if valid {
			next.ServeHTTP(writer, request)
			return
		}
	}
}
