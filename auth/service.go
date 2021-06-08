package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

type Auth struct {
	logger hclog.Logger
}

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("922c5740-d20e-4e52-849c-f13a680bb706"),
})

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(logger hclog.Logger, e *echo.Echo) *Auth {
	service := &Auth{
		logger: logger,
	}

	e.POST("/api/login", service.login)
	e.GET("/api/is-loggedin", service.private, IsLoggedIn)

	return service
}

func (a *Auth) login(c echo.Context) error {
	ld := new(LoginData)
	err := c.Bind(ld)
	if err != nil {
		a.logger.Error("parse login data failed", err)
		return c.String(http.StatusBadRequest, "")
	}

	// Check in your db if the user exists or not
	if ld.Username == "this_is_sparta" && ld.Password == "ww1234@@yagni" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		// This is the information which frontend can use
		// The backend can also decode the token and get admin etc.
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Doe"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Generate encoded token and send it as response.
		// The signing string should be secret (a generated UUID          works too)
		t, err := token.SignedString([]byte("922c5740-d20e-4e52-849c-f13a680bb706"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}

func (a *Auth) private(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
