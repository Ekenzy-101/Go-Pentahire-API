package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/routes"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /auth/login/verify", func() {
	var (
		code         string
		email        string
		OTPSecretKey string
		responseBody gin.H
		token        string
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{"email": email, "code": code}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, "/auth/login/verify", bytes.NewReader(requestBodyBytes))
		if err != nil {
			return nil, err
		}

		request.AddCookie(&http.Cookie{Name: config.VerifyLoginTokenCookieName, Value: token})
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
		token = "y3ryeuyrueiuq"
		email = "test5@test.com"
		responseBody = gin.H{}
	})

	JustBeforeEach(func() {
		key, err := services.GenerateOTPKey(email)
		Expect(err).NotTo(HaveOccurred())

		OTPSecretKey = key.Secret()
		code, err = services.GenerateOTPCode(OTPSecretKey)
		Expect(err).NotTo(HaveOccurred())

		options := models.SQLOptions{
			Arguments:     []interface{}{email, "Test", "Test", "Test", OTPSecretKey},
			InsertColumns: []string{"email", "password", "firstname", "lastname", "otp_secret_key"},
			ReturnColumns: []string{"email"},
			Destination:   []interface{}{&email},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		err = redisClient.Set(ctx, config.RedisVerifyLoginPrefix+email, token, config.RedisVerifyLoginTTL).Err()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_, err := pool.Exec(ctx, "DELETE FROM users")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be a success", func() {
		By("sending a request with valid inputs and token")
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains the user's info")
		actual := helpers.GetMapKeys(responseBody["user"])
		elements := helpers.GetStructFields(models.User{}, []interface{}{"password", "otp_secret_key", "phone_no"})
		Expect(actual).To(ContainElements(elements...))

		By("returning cookies")
		Expect(response.Result().Header).To(HaveKey("Set-Cookie"))
	})

	It("should be an error", func() {
		By("sending a request with invalid inputs")
		email = "invalid email"
		code = "invalid code"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"email", "message"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with an invalid or expired 2fa token")
		token = "invalid token"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 401")
		Expect(response).To(HaveHTTPStatus(http.StatusUnauthorized))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	It("should be an error", func() {
		By("sending a request with an incorrect code")
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

	It("should be an error", func() {
		By("sending a request with an email that does not exist")
		_, err := pool.Exec(ctx, "DELETE FROM users WHERE email =  $1", email)
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 404")
		Expect(response).To(HaveHTTPStatus(http.StatusNotFound))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
