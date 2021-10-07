package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/routes"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /account/otp-key/confirm", func() {
	var (
		accessToken     string
		emailVerifiedAt interface{}
		code            string
		userId          string
		responseBody    gin.H
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{"code": code}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, "/account/otp-key/confirm", bytes.NewReader(requestBodyBytes))
		if err != nil {
			return nil, err
		}

		request.AddCookie(&http.Cookie{Name: config.AccessTokenCookieName, Value: accessToken})
		response := httptest.NewRecorder()
		router := routes.SetupRouter()
		router.ServeHTTP(response, request)
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	BeforeEach(func() {
		responseBody = gin.H{}
		emailVerifiedAt = time.Now()
	})

	JustBeforeEach(func() {
		key, err := services.GenerateOTPKey("Test")
		Expect(err).NotTo(HaveOccurred())

		options := models.SQLOptions{
			Arguments:     []interface{}{"Test", "Test", "Test", "Test", key.Secret(), emailVerifiedAt},
			InsertColumns: []string{"email", "password", "firstname", "lastname", "otp_secret_key", "email_verified_at"},
			ReturnColumns: []string{"id"},
			Destination:   []interface{}{&userId},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		user := &models.User{ID: userId}
		accessToken, err = user.GenerateAccessToken()
		Expect(err).NotTo(HaveOccurred())

		code, err = services.GenerateOTPCode(key.Secret())
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_, err := pool.Exec(ctx, "DELETE FROM users")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be a success", func() {
		By("sending a request with a valid access token")
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains a success message")
		Expect(responseBody).To(HaveKey("message"))
	})

	Context("", func() {
		BeforeEach(func() {
			emailVerifiedAt = nil
		})

		It("should be an error", func() {
			By("sending a request when user's email is not verified")
			response, err := ExecuteRequest()
			Expect(err).NotTo(HaveOccurred())

			By("returning a status code of 400")
			Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

			By("returning a body that contains error messages")
			Expect(responseBody).To(HaveKey("message"))
		})
	})

	It("should be an error", func() {
		By("sending a request with invalid code")
		var err error
		code, err = helpers.GenerateRandomNumbers(6)
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	It("should be a error", func() {
		By("sending a request with an invalid access token")
		token, err := services.SignJWTToken(services.JWTOptions{
			SigningMethod: jwt.SigningMethodHS256,
			Claims: services.AccessTokenClaims{
				ID: userId,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now()),
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		accessToken = token
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 401")
		Expect(response).To(HaveHTTPStatus(http.StatusUnauthorized))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	It("should be an error", func() {
		By("sending a request with a token of a user that doesn't exists in the database")
		_, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userId)
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 404")
		Expect(response).To(HaveHTTPStatus(http.StatusNotFound))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
