package echoSwagger

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/swag"
)

type mockedSwag struct{}

func (s *mockedSwag) ReadDoc() string {
	return `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "Swagger Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {
        "/file/upload": {
            "post": {
                "description": "Upload file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload file",
                "operationId": "file.upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "this is a test file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-string-by-int/{some_id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new pet to the store",
                "operationId": "get-string-by-int",
                "parameters": [
                    {
                        "type": "int",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.Pet"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-struct-array-by-string/{some_id}": {
            "get": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "operationId": "get-struct-array-by-string",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "int",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "int",
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "web.APIError": {
            "type": "object",
            "properties": {
                "CreatedAt": {
                    "type": "string",
                    "format": "date-time"
                },
                "ErrorCode": {
                    "type": "integer"
                },
                "ErrorMessage": {
                    "type": "string"
                }
            }
        },
        "web.Pet": {
            "type": "object",
            "properties": {
                "Category": {
                    "type": "object"
                },
                "ID": {
                    "type": "integer"
                },
                "Name": {
                    "type": "string"
                },
                "PhotoUrls": {
                    "type": "array"
                },
                "Status": {
                    "type": "string"
                },
                "Tags": {
                    "type": "array"
                }
            }
        }
    }
}`
}

func TestWrapHandler(t *testing.T) {
	router := echo.New()

	router.GET("/*", EchoWrapHandler(DocExpansion("none"), DomID("#swagger-ui")))

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 500, w2.Code)

	swag.Register(swag.Name, &mockedSwag{})
	w2 = performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "image/png")

	w4 := performRequest("GET", "/swagger-ui.css", router)
	assert.Equal(t, 200, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "text/css; charset=utf-8")

	w5 := performRequest("GET", "/swagger-ui-bundle.js", router)
	assert.Equal(t, 200, w5.Code)
	assert.Equal(t, w5.Header()["Content-Type"][0], "application/javascript")

	w6 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w6.Code)

	w7 := performRequest("GET", "/", router)
	assert.Equal(t, 301, w7.Code)
}

func TestHandlerReuse(t *testing.T) {
	router := echo.New()

	router.GET("/swagger/*", EchoWrapHandler())
	router.GET("/admin/swagger/*", EchoWrapHandler())

	w1 := performRequest("GET", "/swagger/index.html", router)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w2 := performRequest("GET", "/admin/swagger/index.html", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w3 := performRequest("GET", "/swagger/index.html", router)
	assert.Equal(t, 200, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w4 := performRequest("GET", "/admin/swagger/index.html", router)
	assert.Equal(t, 200, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "text/html; charset=utf-8")
}

func performRequest(method, target string, e *echo.Echo) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	return w
}

func TestURL(t *testing.T) {
	expected := "https://github.com/swaggo/http-swagger"
	cfg := Config{}
	configFunc := URL(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.URL)
}

func TestDeepLinking(t *testing.T) {
	expected := true
	cfg := Config{}
	configFunc := DeepLinking(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DeepLinking)
}

func TestDocExpansion(t *testing.T) {
	expected := "https://github.com/swaggo/docs"
	cfg := Config{}
	configFunc := DocExpansion(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DocExpansion)
}

func TestDomID(t *testing.T) {
	expected := "#swagger-ui"
	cfg := Config{}
	configFunc := DomID(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DomID)
}
