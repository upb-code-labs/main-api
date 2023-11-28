package entities

type MarkdownBlock struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
	Index   int    `json:"index"`
}
