package services

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Pentahire",
		AccountName: email,
	})
}

func ValidateOTP(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
