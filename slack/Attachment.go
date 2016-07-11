package slack

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
  MarkdownIn []string `json:"mrkdwn_in, omitempty"`
}

// Helper type for a slick of Attachments.
type Attachments []Attachment
