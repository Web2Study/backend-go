package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

/*
curl -X POST   -H 'Content-Type: application/json' -d '{"username":"alex","password":"1234"}'  localhost:3000/login

curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"username":"alex","password":"1234"}' \
  localhost:3000/login

curl localhost:1323/restricted -H "Authorization: Bearer {$token}""
*/
func login(c echo.Context) error {
	//name := c.FormValue("name")
	//password := c.FormValue("password")
	u := &User{
		ID: 0,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	log.Println(u.Name, "---", u.Password)
	if u.Password == "1234" { //just for demo ,you should query from DataBase
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = u.Name
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		cookie := new(echo.Cookie)
		cookie.SetName("token")
		cookie.SetValue(t)
		cookie.SetExpires(time.Now().Add(time.Minute * 10))
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
	r := e.Group("/api")
	jwtConfig := middleware.DefaultJWTConfig
	jwtConfig.ContextKey = "user"
	jwtConfig.TokenLookup = "cookie:token"
	jwtConfig.SigningKey = []byte("secret")

	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("", restricted)
	// Routes
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	e.Run(standard.New(":3000"))
}
