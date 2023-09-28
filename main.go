package main

import (
	"vueltop2c/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

// CustomValidator .
type CustomValidator struct {
	validator *validator.Validate
}

// Validate .
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/vueltop2c-mg", routes.ValidarTransaccionVuelto)
	e.POST("/c2pcobro-mg", routes.ValidarTransaccionCobroC2P)
	// e.Logger.Fatal(e.Start(":4446"))
	e.Logger.Fatal(e.StartTLS(":4446", "../cert/snt.crt", "../cert/commercial.key"))
}
