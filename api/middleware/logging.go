package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

//Logging show requests logs
func Logging(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
