package definitions

type BlockRepository interface {
	UpdateMarkdownBlockContent(blockUUID string, content string) (err error)
}
