package slack

type ResponseType int
const (
  // Ephemeral tells Slack the response should be displayed only to the user that
  // executed the slash command. For long running commands, consider immediately
  // responding to the user with an Ephemeral response, and use InChannel for the
  // final result.
  Ephemeral ResponseType = 1 + iota
  // InChannel tells Slack the response should be placed in the channel for all
  // channel members to see.
  InChannel
)
var ResponseTypes = []string {
  "ephemeral",
  "in_hannel",
}
func (r ResponseType) String() string {
  return ResponseTypes[r - 1]
}

// The Response must be json encoded.  Response data must be URL encoded, also.
type Response struct {
  Type string `json:"response_type"`
  Text *string `json:"text"`
  Attachments Attachments `json:"attachments,omitempty"`
}
