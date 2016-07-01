package routes

import (
    "fmt"
    "net/http"

    "github.com/confyrm/gorest/config"
)

// Index handles the '/' path or testing.
func Index(config *config.Config, rw http.ResponseWriter, req *http.Request) error {
    fmt.Fprintln(rw, "Welcome!")
    return nil
}
