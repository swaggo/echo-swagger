package main

import (
	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

func main() {
	e := echo.New()

	e.GET("/swagger/*",echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
