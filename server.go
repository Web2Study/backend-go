package main

import (
	"net/http"

	"time"

	"github.com/GeertJohan/go.rice"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

/*
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"username":"alex","password":"1234"}' \
  localhost:3000/login

curl localhost:1323/restricted -H "Authorization: Bearer {$token}""
*/
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if password == "1234" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		cookie := new(echo.Cookie)
		cookie.SetName("token")
		cookie.SetValue(t)
		cookie.SetExpires(time.Now().Add(24 * time.Hour))
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	e := echo.New()
	// the file server for rice. "app" is the folder where the files come from.
	assetHandler := http.FileServer(rice.MustFindBox("app").HTTPBox())
	// serves the index.html from rice
	e.GET("/", standard.WrapHandler(assetHandler))

	// servers other static files
	e.GET("/static/*", standard.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	//e.File("/favicon.ico", "/static/favicon.ico")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/public", accessible)

	// Restricted group
	r := e.Group("/admin")
	jwtConfig := middleware.DefaultJWTConfig
	jwtConfig.ContextKey = "user"
	jwtConfig.TokenLookup = "cookie:token"
	jwtConfig.SigningKey = []byte("secret")

	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("", restricted)

	e.Run(standard.New(":3000"))
}
