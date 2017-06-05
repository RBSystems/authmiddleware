package authmiddleware

import (
	"errors"
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
		passed, err := checkLocal()
		if err != nil {
			jsonresp.New(writer, http.StatusBadRequest, err.Error())
			return
		}

		if passed {
			next.ServeHTTP(writer, request)
			return
		}

		passed, err = checkBearerToken(request)
		if err != nil {
			jsonresp.New(writer, http.StatusBadRequest, err.Error())
			return
		}

		if passed {
			next.ServeHTTP(writer, request)
			return
		}

		passed, err = checkWSO2(request)
		if err != nil {
			jsonresp.New(writer, http.StatusBadRequest, err.Error())
			return
		}

		if passed {
			next.ServeHTTP(writer, request)
			return
		}

		jsonresp.New(writer, http.StatusBadRequest, "Not authorized")
	})
}

func checkLocal() (bool, error) {
	log.Printf("Local check starting")

	if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
		log.Printf("Authorized via LOCAL_ENVIRONMENT")
		return true, nil
	}

	log.Printf("Local check finished")
	return false, nil
}

func checkBearerToken(request *http.Request) (bool, error) {
	log.Printf("Bearer token check starting")

	token := request.Header.Get("Authorization") // Get the token if it exists

	if len(token) > 0 { // Proceed if we found a token
		parts := strings.Split(token, " ")

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return false, errors.New("Bad Authorization header")
		}

		valid, err := bearertoken.CheckToken([]byte(parts[1])) // Validate the existing token
		if err != nil {
			return false, err
		}

		if valid {
			log.Println("Bearer token authorized")
			return true, nil
		}
	}

	log.Printf("Bearer token check finished")
	return false, nil
}

func checkWSO2(request *http.Request) (bool, error) {
	log.Printf("WSO2 check starting")

	token := request.Header.Get("X-jwt-assertion") // Get the token if it exists

	if len(token) > 0 { // Proceed if we found a token
		valid, err := wso2jwt.Validate(token) // Validate the existing token
		if err != nil {
			log.Printf("Invalid WSO2 information")
			return false, err
		}

		if valid {
			log.Printf("WSO2 validated successfully")
			return true, nil
		}
	}

	log.Printf("WSO2 check finished")
	return false, nil
}
