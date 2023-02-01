package echoSwagger

import (
	"io/ioutil"
	"net/http"
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

	router.Any("/*", EchoWrapHandler(DocExpansion("none"), DomID("swagger-ui")))

	w1 := performRequest(http.MethodGet, "/index.html", router)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	assert.Equal(t, http.StatusInternalServerError, performRequest(http.MethodGet, "/doc.json", router).Code)

	doc := &mockedSwag{}
	swag.Register(swag.Name, doc)
	w2 := performRequest(http.MethodGet, "/doc.json", router)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "application/json; charset=utf-8")

	// Perform body rendering validation
	w2Body, err := ioutil.ReadAll(w2.Body)
	assert.NoError(t, err)
	assert.Equal(t, doc.ReadDoc(), string(w2Body))

	w3 := performRequest(http.MethodGet, "/favicon-16x16.png", router)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "image/png")

	w4 := performRequest(http.MethodGet, "/swagger-ui.css", router)
	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "text/css; charset=utf-8")

	w5 := performRequest(http.MethodGet, "/swagger-ui-bundle.js", router)
	assert.Equal(t, http.StatusOK, w5.Code)
	assert.Equal(t, w5.Header()["Content-Type"][0], "application/javascript")

	assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/notfound", router).Code)

	assert.Equal(t, http.StatusMovedPermanently, performRequest(http.MethodGet, "/", router).Code)

	assert.Equal(t, http.StatusMethodNotAllowed, performRequest(http.MethodPost, "/index.html", router).Code)

	assert.Equal(t, http.StatusMethodNotAllowed, performRequest(http.MethodPut, "/index.html", router).Code)

}

func TestConfig(t *testing.T) {
	router := echo.New()

	swaggerHandler := URL("swagger.json")
	router.Any("/*", EchoWrapHandler(swaggerHandler))

	w := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `url: "swagger.json"`)
}

func TestConfigWithOAuth(t *testing.T) {
	router := echo.New()

	swaggerHandler := EchoWrapHandler(OAuth(&OAuthConfig{
		ClientId: "my-client-id",
		Realm:    "my-realm",
		AppName:  "My App Name",
	}))
	router.GET("/*", swaggerHandler)

	w := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `ui.initOAuth({
    clientId: "my-client-id",
    realm: "my-realm",
    appName: "My App Name",
    usePkceWithAuthorizationCodeGrant:  false 
  })`)
}

func TestHandlerReuse(t *testing.T) {
	router := echo.New()

	router.GET("/swagger/*", EchoWrapHandler())
	router.GET("/admin/swagger/*", EchoWrapHandler())

	w1 := performRequest(http.MethodGet, "/swagger/index.html", router)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w2 := performRequest(http.MethodGet, "/admin/swagger/index.html", router)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w3 := performRequest(http.MethodGet, "/swagger/index.html", router)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w4 := performRequest(http.MethodGet, "/admin/swagger/index.html", router)
	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "text/html; charset=utf-8")
}

type httpWriter struct{}

func (h httpWriter) Header() http.Header {
	return http.Header{}
}

func (h httpWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (h httpWriter) WriteHeader(_ int) {}

func TestMissingFlusher(t *testing.T) {
	router := echo.New()

	router.GET("/swagger/*", EchoWrapHandler())

	r := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	router.ServeHTTP(httpWriter{}, r)
}

func performRequest(method, target string, e http.Handler) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	return w
}

func TestURL(t *testing.T) {
	var cfg Config
	expected := "https://github.com/swaggo/http-swagger"
	URL(expected)(&cfg)
	assert.Equal(t, expected, cfg.URL)
}

func TestDeepLinking(t *testing.T) {
	var cfg Config
	expected := true
	DeepLinking(expected)(&cfg)
	assert.Equal(t, expected, cfg.DeepLinking)
}

func TestSyntaxHighlight(t *testing.T) {
	var cfg Config
	expected := true
	SyntaxHighlight(expected)(&cfg)
	assert.Equal(t, expected, cfg.SyntaxHighlight)
}

func TestDocExpansion(t *testing.T) {
	var cfg Config
	expected := "https://github.com/swaggo/docs"
	DocExpansion(expected)(&cfg)
	assert.Equal(t, expected, cfg.DocExpansion)
}

func TestDomID(t *testing.T) {
	var cfg Config
	expected := "swagger-ui"
	DomID(expected)(&cfg)
	assert.Equal(t, expected, cfg.DomID)
}

func TestInstanceName(t *testing.T) {
	var cfg Config

	expected := "custom-instance-name"
	InstanceName(expected)(&cfg)
	assert.Equal(t, expected, cfg.InstanceName)

	newCfg := newConfig(InstanceName(""))
	assert.Equal(t, swag.Name, newCfg.InstanceName)
}

func TestPersistAuthorization(t *testing.T) {
	var cfg Config
	expected := true
	PersistAuthorization(expected)(&cfg)
	assert.Equal(t, expected, cfg.PersistAuthorization)
}

func TestOAuth(t *testing.T) {
	var cfg Config
	expected := OAuthConfig{
		ClientId: "my-client-id",
		Realm:    "my-realm",
		AppName:  "My App Name",
		UsePkce:  true,
	}
	OAuth(&expected)(&cfg)
	assert.Equal(t, expected.ClientId, cfg.OAuth.ClientId)
	assert.Equal(t, expected.Realm, cfg.OAuth.Realm)
	assert.Equal(t, expected.AppName, cfg.OAuth.AppName)
	assert.Equal(t, expected.UsePkce, cfg.OAuth.UsePkce)
}

func TestOAuthNil(t *testing.T) {
	var cfg Config
	var expected *OAuthConfig
	OAuth(expected)(&cfg)
	assert.Equal(t, expected, cfg.OAuth)
}
