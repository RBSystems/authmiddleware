# wso2jwt

JSON wants returned values to be in key-value format. This package removes a couple lines of code I got tired of copypasta-ing between projects.

## Dependencies
- [Echo](https://labstack.com/echo)

## Installation
```
go get github.com/byuoitav/wso2jwt
```

## Usage
```
import "github.com/byuoitav/wso2jwt"
```
```
func main() {
  port := ":8000"
  router := echo.New()
  router.Use(wso2jwt.ValidateJWT())

  router.Run(fasthttp.New(port))
}
```
