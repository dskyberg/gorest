package routes

import (
    "os"
    "fmt"
    "net/http"

    "github.com/confyrm/gorest/config"
)

func Exit(config *config.Config, rw http.ResponseWriter, req *http.Request) error {
    fmt.Fprintln(rw, "Exiting")
    os.Exit(0)
    return nil
}
