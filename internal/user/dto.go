package user

type StoreDto struct {
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Role     string `json:"role_id" form:"role_id" validate:"required"`
	Image    string `json:"img_temp" form:"img_temp"`
}
type UpdateDto struct {
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password"`
	Role     string `json:"role_id" form:"role_id" form:"role_id" validate:"required"`
}
