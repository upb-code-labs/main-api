package entities

type Language struct {
	UUID                string `json:"uuid"`
	TemplateArchiveUUID string `json:"-"`
	Name                string `json:"name"`
}
