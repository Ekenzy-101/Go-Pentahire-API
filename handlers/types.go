package handlers

type CodeField struct {
	Code string `json:"code" binding:"required,len=6"`
}

type EmailField struct {
	Email string `json:"email" binding:"email,max=255"`
}

type NameFields struct {
	Firstname string `json:"firstname" binding:"required,name,max=50"`
	Lastname  string `json:"lastname" binding:"required,name,max=50"`
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
	NameFields
	TokenField
	LoginRequestBody
}

type ResetPasswordRequestBody struct {
	TokenField
	PasswordField
}

type UpdateProfileRequestBody struct {
	NameFields
	EmailField
}

type UpdatePasswordRequestBody struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=128,password"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128,password"`
}

type VerifyLoginRequestBody struct {
	EmailField
	CodeField
}
