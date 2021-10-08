package config

import "time"

const (
	AccessTokenCookieName        = "pnt_acc_token"
	AccessTokenTTLInSeconds      = 60 * 60
	RefreshTokenCookieName       = "pnt_ref_token"
	RefreshTokenTTLInSeconds     = 60 * 60 * 24 * 7
	VerifyLoginTokenCookieName   = "pnt_2fa_token"
	VerifyLoginTokenTTLInSeconds = 60 * 5

	RedisResetPasswordPrefix = "reset_password:"
	RedisResetPasswordTTL    = 1 * time.Hour
	RedisVerifyEmailPrefix   = "verify_email:"
	RedisVerifyEmailTTL      = 24 * time.Hour
	RedisVerifyLoginPrefix   = "verify_login:"
	RedisVerifyLoginTTL      = 5 * time.Minute

	UsersTable = "users"
)
