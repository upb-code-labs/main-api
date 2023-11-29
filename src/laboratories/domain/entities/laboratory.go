package entities

type Laboratory struct {
	UUID           string          `json:"uuid"`
	CourseUUID     string          `json:"-"`
	RubricUUID     *string         `json:"rubric_uuid"`
	Name           string          `json:"name"`
	OpeningDate    string          `json:"opening_date"`
	DueDate        string          `json:"due_date"`
	MarkdownBlocks []MarkdownBlock `json:"markdown_blocks"`
	TestBlocks     []TestBlock     `json:"test_blocks"`
}
