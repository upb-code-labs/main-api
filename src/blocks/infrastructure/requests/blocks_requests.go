package requests

type UpdateMarkdownBlockContentRequest struct {
	Content string `json:"content" validate:"required"`
}
