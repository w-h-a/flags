package slack

type slackMessage struct {
	Text         string       `json:"text"`
	Attachements []attachment `json:"attachments"`
}

type attachment struct {
	Color  string  `json:"color"`
	Title  string  `json:"title"`
	Fields []field `json:"fields"`
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
}
