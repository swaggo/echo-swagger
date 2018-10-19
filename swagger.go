package echoSwagger

import (
	"golang.org/x/net/webdav"
	"html/template"
	"net/http"
	"regexp"

	"github.com/labstack/echo"
	"github.com/swaggo/files"
	"github.com/swaggo/swag"
)

// WrapHandler wraps swaggerFiles.Handler and returns echo.HandlerFunc
var WrapHandler = wrapHandler(swaggerFiles.Handler)

// wapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func wrapHandler(h *webdav.Handler) echo.HandlerFunc {
	//create a template with name
	t := template.New("swagger_index.html")
	index, _ := t.Parse(indexTempl)

	type pro struct {
		Host string
	}

	var re = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[\?|.]*`)

	return func(c echo.Context) error {
		var matches []string
		if matches = re.FindStringSubmatch(c.Request().RequestURI); len(matches) != 3 {

			return c.String(http.StatusNotFound, "404 page not found")
		}
		path := matches[2]
		prefix := matches[1]
		h.Prefix = prefix

		switch path {
		case "index.html":
			s := &pro{
				Host: "doc.json", //TODO: provide to customs?
			}
			index.Execute(c.Response().Writer, s)
		case "doc.json":
			doc, _ := swag.ReadDoc()
			c.Response().Write([]byte(doc))
		default:
			h.ServeHTTP(c.Response().Writer, c.Request())

		}

		return nil
	}
}

const indexTempl = `<!DOCTYPE html>
<html>
  <head>
    <title>ReDoc</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

    <!--
    ReDoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='{{.Host}}'></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`
