package auth


type LoginDto struct {
	Password string `json:"password" form:"password" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
}

