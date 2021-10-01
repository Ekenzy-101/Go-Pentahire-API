package config

const (
	AccessTokenCookieName    = "pnt_acc_token"
	AccessTokenTTLInSeconds  = 60 * 60
	RefreshTokenCookieName   = "pnt_ref_token"
	RefreshTokenTTLInSeconds = 60 * 60 * 24 * 7

	RedisVerifyEmailPrefix   = "verify_email:"
	RedisResetPasswordPrefix = "reset_password:"

	UsersTable = "users"
)
