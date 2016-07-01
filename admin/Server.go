package admin

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
    log.Printf("Admin: Listening on %d...", config.GetInt("ADMIN_PORT"))
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetInt("ADMIN_PORT")), loggedRouter))
}
