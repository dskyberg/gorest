package app

import (
    "os"
    "log"
    "fmt"
    "net/http"

    "github.com/gorilla/handlers"
    "github.com/confyrm/gorest/config"
    . "github.com/confyrm/gorest/admin/routes"
)

func Server(config *config.Config)  {

    router := RouteSet.New(config)
    loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)
    log.Printf("%s: Listening on %d...", config.GetString("APP_NAME"), config.GetInt("APP_PORT"))
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetInt("APP_PORT")), loggedRouter))
}
