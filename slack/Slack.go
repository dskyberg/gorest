// Package slack provides a set of types for using the Slack API. This includes
// types for custom commands to receive REST requests from Slack.
package slack

// The Slack Request that is sent by Slack to a custom app.
type Request struct {
  // The token that was used by Slack.  This should be validated against
  // SLACK_TOKEN.
  Token string `schema:"token"`
  // The Slack team the below user is part of.
  TeamId string `schema:"team_id"`
  // The Slack team domain.
  TeamDomain string `schema:"team_domain"`
  // The Slack channel the slash command was executed on.
  ChannelId string `schema:"channel_id"`
  // The name of the Slack channel the slash command was executed on.
  ChannelName string `schema:"channel_name"`
  // The Id of the user executing the slash command.
  UserId string `schema:"user_id"`
  // The Slack user name of the user that executed the slash command.
  UserName string `schema:"user_name"`
  // The Slack slash command that was executed.
  Command string `schema:"command"`
  // all text that was provided after the slash command.
  Text string `schema:"text"`
  // If the time to respond to the command is long, the response can be
  // sent to this url, instead ofo in the initial response.
  ResponseUrl string `schema:"response_url"`
}

// InChannel tells Slack the response should be placed in the channel for all
// channel members to see.
const InChannel = "in_channel"
// Ephemeral tells Slack the response should be displayed only to the user that
// executed the slash command. For long running commands, consider immediately
// responding to the user with an Ephemeral response, and use InChannel for the
// final result.
const Ephemeral = "ephemeral"

// Slash standard green color for attachments.
const GOOD = "good"
// Slash standard yellow color for attachments.
const WARNING = "warning"
// Slash standard red color for attachments.
const DANGER = "danger"

// AttachmentFields are displayed as tabular data.
type AttachmentField struct {
  // Table header value.
  Title string `json:"title"`
  // Table value.  Displayed under the header.
  Value string `json:"value"`
  // true, if this field is short enough to be displayed on same line as
  // previous fields.  Use false, if you know it is too long
  Short bool `json:"short"`
}
// Helper type for a slice of AttachmentFields.
type AttachmentFields []AttachmentField

// Attachment is the primary tool for creating rich responses.
type Attachment struct {
  Title string `json:"title"`
  TitleLink string `json:"title_link, omitempty"`

  PreText string `json:"pretext, omitempty"`
  Text string `json:"text, omitempty"`
  Color string `json:"color, omitempty"`
  Fallback string `json:"fallback"`

  ImageUrl string `json:"image_url, omitempty"`
  ThumbUrl string `json:"thumb_url, omitempty"`

  // Author fields
  AuthorName string `json:"author_name, omitempty"`
  AuthorLink string `json:"author_link, omitempty"`
  AuthorIcon string `json:"author_icon, omitempty"`


  // Footer fields
  Footer string `json:"footer, omitempty"`
  FooterIcon string `json:"footer_icon, omitempty"`
  TimeStamp int `json:"ts", omitempty`

  // Table fields
  Fields AttachmentFields `json:"fields, omitempty"`
}

// Helper type for a slick of Attachments.
type Attachments []Attachment

// The Response must be json encoded.  Response data must be URL encoded, also.
type Response struct {
  ResponseType string `json:"response_type"`
  Text string `json:"text"`
  Attachments Attachments `json:"attachments,omitempty"`
}
