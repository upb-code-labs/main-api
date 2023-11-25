package entities

type MarkdownBlock struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}
