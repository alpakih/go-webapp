package permission

type StoreDto struct {
	Group       string `json:"group" form:"group" validate:"required"`
	Feature     string `json:"feature" form:"feature" validate:"required"`
	Url         string `json:"url" form:"url" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type UpdateDto struct {
	Group       string `json:"group" form:"group" validate:"required"`
	Feature     string `json:"feature" form:"feature" validate:"required"`
	Url         string `json:"url" form:"url" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}
