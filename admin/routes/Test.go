package routes

import (
    "log"
    "io/ioutil"
    "net/http"
    //"encoding/json"

    "github.com/confyrm/gorest/config"
    . "github.com/confyrm/gorest/errors"
)

// Test is just a dumb POST handler that logs whatever was posted.
func Test(config *config.Config, rw http.ResponseWriter, req *http.Request) error {

  // If the post contains structured json, decode it into a struct like this:
  /*
  decoder := json.NewDecoder(req.Body)
  // Change test_struct to your type
  var t test_struct
  if err := decoder.Decode(&t); err != nil {
    // Handle the decoding error.
  }
  */

  // We are just grabbing the body as a string.
  body, err := ioutil.ReadAll(req.Body);
  if err != nil {
    log.Printf("Admin.Test Error reading request body: %#v\n\n", err)
    return StatusError{http.StatusInternalServerError, err}
  }
  // Just dump the request body to the log.
  log.Printf("Admin.Test received: %s\n\n", body)
  return nil
}
