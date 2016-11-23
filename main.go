package authmiddleware

import (
	"net/http"
	"os"

	"github.com/jessemillar/jsonresp"
)

// Authenticate is the middleware function
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if len(os.Getenv("LOCAL_ENVIRONMENT")) == 0 { // If the `LOCAL_ENVIRONMENT` environment variable isn't set, proceed
			token := request.Header.Get("X-jwt-assertion")
			if token == "" {
				jsonresp.New(writer, http.StatusBadRequest, "No WSO2-provided `X-jwt-assertion` header present")
				return
			}

			err := wso2jwt.Validate(token)
			if err != nil {
				jsonresp.New(writer, http.StatusBadRequest, err.Error())
				return
			}
		}

		next.ServeHTTP(writer, request)
	})
}
