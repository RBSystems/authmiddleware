package authmiddleware

import (
	"log"
	"net/http"
	"os"

	"github.com/byuoitav/authmiddleware/wso2jwt"
	"github.com/jessemillar/jsonresp"
)

// Authenticate is the middleware function
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// First, check if the `LOCAL_ENVIRONMENT` environment variable is set and if it is, return without further checks (circumventing the authentication for local development and Raspberry Pi's)
		log.Printf("Local check starting")
		if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
			log.Printf("Local check finished")
			next.ServeHTTP(writer, request)
		}

		log.Printf("WSO2 check starting")
		// Next, check for the presence of a JWT token provided by WSO2, if there is one, go through the process of validating it and then return true if it's valid
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
	})
}
