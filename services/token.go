package services

import "github.com/golang-jwt/jwt/v4"

type JWTOption struct {
	jwt.SigningMethod
	Claims jwt.Claims
	Secret string
	Token  string
}

type AccessTokenClaim struct {
	Email string `json:"email"`
	ID    string `json:"_id"`
	jwt.RegisteredClaims
}

func SignToken(option JWTOption) (string, error) {
	token := jwt.NewWithClaims(option.SigningMethod, option.Claims)
	signedToken, err := token.SignedString([]byte(option.Secret))

	return signedToken, err
}

func VerifyToken(option JWTOption) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		option.Token,
		option.Claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(option.Secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
