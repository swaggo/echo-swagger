package echoSwagger

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	_ "github.com/swaggo/gin-swagger/example/docs"
)

func TestWrapHandler(t *testing.T) {

	router := echo.New()

	router.GET("/*", WrapHandler)

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, "text/html; charset=utf-8", w1.Header().Get("Content-Type"))

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "application/json", w2.Header().Get("Content-Type"))

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)
	assert.Equal(t, "image/png", w3.Header().Get("Content-Type"))

	w4 := performRequest("GET", "/swagger-ui-bundle.js", router)
	assert.Equal(t, 200, w4.Code)
	assert.Equal(t, "application/javascript", w4.Header().Get("Content-Type"))

	w5 := performRequest("GET", "/swagger-ui.css", router)
	assert.Equal(t, 200, w5.Code)
	assert.Equal(t, "text/css; charset=utf-8", w5.Header().Get("Content-Type"))

	w6 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w6.Code)
}

func performRequest(method, target string, e *echo.Echo) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	return w
}
