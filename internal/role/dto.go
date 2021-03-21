package role

type Store struct {
	Slug        string   `json:"slug" form:"slug" validate:"required"`
	RoleName    string   `json:"role_name" form:"role_name" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	Permission  []string `json:"permissions" form:"permissions[]"`
}

type Update struct {
	Slug        string   `json:"slug" form:"slug" validate:"required"`
	RoleName    string   `json:"role_name" form:"role_name" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	Permission  []string `json:"permissions" form:"permissions[]"`
}
