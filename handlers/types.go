package handlers

type CodeField struct {
	Code string `json:"code" binding:"required,max=6"`
}

type EmailField struct {
	Email string `json:"email" binding:"email,max=255"`
}

type PasswordField struct {
	Password string `json:"password" binding:"required,min=8,max=128,password"`
}

type TokenField struct {
	Token string `json:"token" binding:"required"`
}

type LoginRequestBody struct {
	EmailField
	PasswordField
}

type RegisterRequestBody struct {
	Firstname string `json:"firstname" binding:"required,name,max=50"`
	Lastname  string `json:"lastname" binding:"required,name,max=50"`
	TokenField
	LoginRequestBody
}

type ResetPasswordRequestBody struct {
	TokenField
	PasswordField
}
