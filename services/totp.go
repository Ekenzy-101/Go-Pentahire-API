package services

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOTPKey(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Pentahire",
		AccountName: email,
	})
}

func ValidateOTP(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}

func GenerateOTPCode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
