package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var (
	ApiResponseForbidden        = SimpleResponse{Status: http.StatusForbidden, Message: "Seems like you are not welcome here... Bye"}
	ApiResponseBadJson          = SimpleResponse{Status: http.StatusBadRequest, Message: "Expect valid json payload"}
	ApiResponseDuplicatedId     = SimpleResponse{Status: http.StatusConflict, Message: "Duplicated Id"}
	ApiResponseBadRSAEncryption = SimpleResponse{Status: http.StatusUnprocessableEntity, Message: "Error validating payload, bad public key"}
	ApiResponseNoPayload        = SimpleResponse{Status: http.StatusUnprocessableEntity, Message: "No payload"}
	ApiResponseBadRequest       = SimpleResponse{Status: http.StatusBadRequest, Message: "Bad Request"}
	ApiResponseResourceNotFound = SimpleResponse{Status: http.StatusTeapot, Message: "Resource Not Found"}
	ApiResponseNotFound         = SimpleResponse{Status: http.StatusNotFound, Message: "Not Found"}

	// RSA Private key
	// Automatically injected on autobuild with make
	PRIV_KEY = []byte(`INJECT_PRIV_KEY_HERE`)

	// BuntDB Database for store the keys
	// It will create if not exists
	Database = "./database.db"
)

type SimpleResponse struct {
	Status  int
	Message string
}

func main() {
	// Start the server
	e := echo.New()
	e.SetHTTPErrorHandler(CustomHTTPErrorHandler)

	e.Use(middleware.CORS())

	e.POST("/api/keys/add", addKeys)
	e.GET("/api/keys/:id", getEncryptionKey)

	log.Println("Listening on port 8080")
	e.Run(standard.New(":8080"))
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		// If is an API call return a JSON response
		path := c.Request().URL().Path()
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}

		if strings.Contains(path, "/api/") {
			c.JSON(httpError.Code, SimpleResponse{Status: httpError.Code, Message: httpError.Message})
			return
		}

		// Otherwise return the normal response
		c.String(httpError.Code, httpError.Message)
	}
}
