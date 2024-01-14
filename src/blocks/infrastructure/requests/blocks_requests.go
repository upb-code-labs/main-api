package requests

type UpdateMarkdownBlockContentRequest struct {
	Content string `json:"content" validate:"required"`
}

type UpdateTestBlockRequest struct {
	LanguageUUID string `validate:"required,uuid4"`
	Name         string `validate:"required,min=4,max=255"`
}

type SwapBlocksRequest struct {
	FirstBlockUUID  string `json:"first_block_uuid" validate:"required,uuid4"`
	SecondBlockUUID string `json:"second_block_uuid" validate:"required,uuid4"`
}
