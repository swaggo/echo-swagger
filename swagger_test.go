package echoSwagger

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	_ "github.com/swaggo/gin-swagger/example/docs"
)

func TestWrapHandler(t *testing.T) {
	for _, prefix := range []string{"/", "/swagger/"} {
		router := echo.New()
		RegisterEchoHandler(router.GET, prefix, nil)

		w1 := performRequest("GET", prefix+"index.html", router)
		assert.Equal(t, 200, w1.Code, prefix)
		indexContent := w1.Body.Bytes()

		w1_2 := performRequest("GET", prefix+"", router)
		assert.Equal(t, 200, w1_2.Code, prefix)
		assert.Equal(t, indexContent, w1_2.Body.Bytes())

		w2 := performRequest("GET", prefix+"doc.json", router)
		assert.Equal(t, 200, w2.Code, prefix)

		w3 := performRequest("GET", prefix+"favicon-16x16.png", router)
		assert.Equal(t, 200, w3.Code, prefix)

		w3_1 := performRequest("GET", prefix+"swagger-ui.css", router)
		assert.Equal(t, 200, w3_1.Code, prefix)

		w4 := performRequest("GET", prefix+"notfound", router)
		assert.Equal(t, 404, w4.Code, prefix)
	}
}

func performRequest(method, target string, e *echo.Echo) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	return w
}
