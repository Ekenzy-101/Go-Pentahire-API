package config

import "time"

const (
	AccessTokenCookieName    = "pnt_acc_token"
	AccessTokenTTLInSeconds  = 60 * 60
	RefreshTokenCookieName   = "pnt_ref_token"
	RefreshTokenTTLInSeconds = 60 * 60 * 24 * 7

	RedisVerifyEmailPrefix   = "verify_email:"
	RedisVerifyEmailTTL      = 24 * time.Hour
	RedisResetPasswordPrefix = "reset_password:"
	RedisResetPasswordTTL    = 1 * time.Hour

	UsersTable = "users"
)
