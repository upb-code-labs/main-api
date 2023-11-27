package entities

type Laboratory struct {
	UUID           string          `json:"uuid"`
	CourseUUID     string          `json:"course_uuid"`
	RubricUUID     string          `json:"rubric_uuid"`
	Name           string          `json:"name"`
	OpeningDate    string          `json:"opening_date"`
	DueDate        string          `json:"due_date"`
	MarkdownBlocks []MarkdownBlock `json:"markdown_blocks"`
	TestBlocks     []TestBlock     `json:"test_blocks"`
}
