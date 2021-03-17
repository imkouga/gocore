package route

import (
	"net/http"

	"github.com/imkouga/gocore/http/httpserver"
)

func RawHandleFunc(path string, handler func(w http.ResponseWriter, r *http.Request)) error {
	return httpserver.HandleFunc(path, handler)
}

func RawHandlePrefixPathFunc(prefix string, handler func(w http.ResponseWriter, r *http.Request)) error {
	return httpserver.HandlePrefixPathFunc(prefix, handler)
}
