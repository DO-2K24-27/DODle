package apisecurity

import (
	"net/http"
	"os"
)

func IsAuthorized(r *http.Request) bool {
	return r.Header.Get("API-Token") == os.Getenv("API_TOKEN")
}
