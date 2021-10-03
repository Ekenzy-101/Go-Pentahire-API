package services

import "github.com/golang-jwt/jwt/v4"

type JWTOptions struct {
	jwt.SigningMethod
	Claims jwt.Claims
	Secret string
	Token  string
}

type AccessTokenClaims struct {
	Email string `json:"email"`
	ID    string `json:"_id"`
	jwt.RegisteredClaims
}

func SignJWTToken(options JWTOptions) (string, error) {
	token := jwt.NewWithClaims(options.SigningMethod, options.Claims)
	signedToken, err := token.SignedString([]byte(options.Secret))

	return signedToken, err
}

func VerifyJWTToken(options JWTOptions) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		options.Token,
		options.Claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(options.Secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
