# wso2jwt

Validating [JWT tokens](https://jwt.io/) with BYU's implementation of [WSO2](http://wso2.com/products/api-manager/) isn't the easiest thing on the planet. We wanted to make sure that we protected every route with a JWT check and, as part of that check, pulled in the latest signing key.

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

## Notes
If you're not interested in protecting every single route (putting [`/health`](https://github.com/jessemillar/health) behind authentication just seems silly), you can use [Echo Groups](https://echo.labstack.com/middleware/overview) instead of `router.Use()`:
```
router.Group("/command", wso2jwt.ValidateJWT())
```
