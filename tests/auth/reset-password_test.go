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
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /auth/reset-password", func() {
	var (
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
		responseBody = gin.H{}
		userId = uuid.NewString()
	})

	JustBeforeEach(func() {
		options := models.SQLOptions{
			Arguments:     []interface{}{"test2@test.com", "testing", "Test", "Test"},
			InsertColumns: []string{"email", "password", "firstname", "lastname"},
			ReturnColumns: []string{"id"},
			Destination:   []interface{}{&userId},
		}
		response := models.InsertUserRow(ctx, options)
		Expect(response).To(BeNil())

		var err error
		token, err = helpers.GenerateRandomToken(24)
		Expect(err).NotTo(HaveOccurred())

		err = redisClient.Set(ctx, config.RedisResetPasswordPrefix+token, userId, config.RedisResetPasswordTTL).Err()
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

		By("clearing token key from redis")
		err = redisClient.Get(ctx, config.RedisResetPasswordPrefix+token).Err()
		Expect(err).To(MatchError(redis.Nil))

		By("returning a body that contains a success message")
		Expect(responseBody).To(HaveKey("message"))
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
		elements := []interface{}{"password", "message"}
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
		err := redisClient.Set(ctx, config.RedisResetPasswordPrefix+token, uuid.NewString(), config.RedisResetPasswordTTL).Err()
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 404")
		Expect(response).To(HaveHTTPStatus(http.StatusNotFound))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
