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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /auth/reset-password", func() {
	var (
		OTPSecretKey string
		password     string
		token        string
		userId       string
		responseBody gin.H
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{"token": token, "password": password}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewReader(requestBodyBytes))
		if err != nil {
			return nil, err
		}

		response := httptest.NewRecorder()
		router := routes.SetupRouter()
		router.ServeHTTP(response, request)
		if responseBody != nil {
			err = json.NewDecoder(response.Body).Decode(&responseBody)
			if err != nil {
				return nil, err
			}
		}

		return response, nil
	}

	BeforeEach(func() {
		password = "Testing@123"
		OTPSecretKey = ""
		responseBody = gin.H{}
		userId = uuid.NewString()

		var err error
		token, err = helpers.GenerateRandomToken(24)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		options := models.SQLOptions{
			Arguments:     []interface{}{"test2@test.com", "testing", "Test", "Test", OTPSecretKey},
			InsertColumns: []string{"email", "password", "firstname", "lastname", "otp_secret_key"},
			ReturnColumns: []string{"id"},
			Destination:   []interface{}{&userId},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		err := redisClient.Set(ctx, config.RedisResetPasswordPrefix+token, userId, 1*time.Hour).Err()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_, err := pool.Exec(ctx, "DELETE FROM users")
		Expect(err).NotTo(HaveOccurred())

		err = redisClient.FlushDBAsync(ctx).Err()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be a success", func() {
		By("sending a request with valid inputs")
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains the user's info if user's 2FA is disabled")
		actual := helpers.GetMapKeys(responseBody["user"])
		elements := helpers.GetStructFields(models.User{}, []interface{}{"password", "otp_secret_key"})
		Expect(actual).To(ContainElements(elements...))

		By("returning cookies if user's 2FA is disabled")
		Expect(response.Result().Header).To(HaveKey("Set-Cookie"))
	})

	Context("", func() {
		BeforeEach(func() {
			OTPSecretKey = "testsecret"
		})

		It("should be a success", func() {
			By("sending a request with valid inputs")
			responseBody = nil
			response, err := ExecuteRequest()
			Expect(err).NotTo(HaveOccurred())

			By("returning a status code of 204")
			Expect(response).To(HaveHTTPStatus(http.StatusNoContent))

			By("returning an empty body if user's 2FA is enabled")
			Expect(response).To(HaveHTTPBody([]byte(nil)))

			By("not returning cookies if user's 2FA is enabled")
			Expect(response.Result().Header).NotTo(HaveKey("Set-Cookie"))
		})
	})

	It("should be an error", func() {
		By("sending a request with invalid inputs")
		password = "invalid password"
		token = ""
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"password", "token"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with a token that does not exist or has expired")
		token = "does not exist"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	It("should be an error", func() {
		By("sending a request when user is not found")
		err := redisClient.Set(ctx, config.RedisResetPasswordPrefix+token, uuid.NewString(), 1*time.Hour).Err()
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 404")
		Expect(response).To(HaveHTTPStatus(http.StatusNotFound))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
