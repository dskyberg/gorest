package routes

import (
    "os"
    "fmt"
    "net/http"

    "github.com/confyrm/gorest/config"
)
// My first admin method.  This just causes the server to shud down in a
// somewhat graceful fashion.
func Exit(config *config.Config, rw http.ResponseWriter, req *http.Request) error {
    fmt.Fprintln(rw, "Exiting")
    os.Exit(0)
    return nil
}
