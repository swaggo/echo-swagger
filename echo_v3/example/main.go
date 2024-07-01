package main

import (
	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger/echo_v3"
	_ "github.com/swaggo/echo-swagger/example/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Swagger Example API With Echo V3
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	e := echo.New()

	e.GET("/swagger/*", echo_v3.WrapHandler)

	/*
		Or can use EchoWrapHandler func with configurations.
		url := echoSwagger.URL("http://localhost:1323/swagger/doc.json") //The url pointing to API definition
		e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	*/
	e.Logger.Fatal(e.Start(":1323"))
}
