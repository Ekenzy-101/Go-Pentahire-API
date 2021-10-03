package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
)

type HCaptchaResponseBody struct {
	ChallengeTimestamp time.Time `json:"challenge_ts"`
	Credit             bool      `json:"credit"`
	ErrorCodes         []string  `json:"error-codes"`
	Hostname           string    `json:"hostname"`
	Success            bool      `json:"success"`
}

func VerifyHCaptchaToken(ctx context.Context, token string) (*HCaptchaResponseBody, error) {
	const verifyURL = "https://hcaptcha.com/siteverify"
	formData := url.Values{}
	formData.Set("secret", config.HCaptchaSecretKey)
	formData.Set("response", token)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBody := &HCaptchaResponseBody{}
	err = json.NewDecoder(response.Body).Decode(responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
