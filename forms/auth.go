package forms

type AuthLoginForm struct {
	UserName string `json:"user_name" validate:"required,min=3,max=10"`
	Passwd   string `json:"passwd" validate:"required,min=3,max=10"`
}
