package handler

import (
  "log"
  "net/http"
  "github.com/confyrm/gorest/config"
  . "github.com/confyrm/gorest/errors"
)

// EnvHandlerFunc is just like a standard HttpHandler, but with an added
// Config param. This allows us to pass the Config down from main to
// each handler.
type EnvHandlerFunc func(e *config.Config, rw http.ResponseWriter, req *http.Request) error

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
  Env *config.Config
  H EnvHandlerFunc
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  err := h.H(h.Env, rw, req)
  if err != nil {
    switch e := err.(type) {
    case Error:
      // We can retrieve the status here and write out a specific
      // HTTP status code.
      log.Printf("HTTP %d - %s", e.Status(), e)
      http.Error(rw, e.Error(), e.Status())
    default:
      // Any error types we don't specifically look out for default
      // to serving a HTTP 500
      http.Error(rw, http.StatusText(http.StatusInternalServerError),
        http.StatusInternalServerError)
    }
  }
}
