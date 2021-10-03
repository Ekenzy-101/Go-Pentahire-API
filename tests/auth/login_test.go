package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/routes"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /auth/login", func() {
	var (
		email        string
		OTPSecretKey string
		password     string
		responseBody gin.H
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{"email": email, "password": password}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(requestBodyBytes))
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
		email = "test1@test.com"
		OTPSecretKey = ""
		password = "Testing@123"
		responseBody = gin.H{}
	})

	JustBeforeEach(func() {
		user := &models.User{Email: email, Password: password}
		Expect(user.HashPassword()).To(Succeed())

		options := models.SQLOptions{
			Arguments:     []interface{}{user.Email, user.Password, "Test", "Test", OTPSecretKey},
			InsertColumns: []string{"email", "password", "firstname", "lastname", "otp_secret_key"},
			ReturnColumns: []string{"email"},
			Destination:   []interface{}{&user.Email},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())
	})

	AfterEach(func() {
		_, err := pool.Exec(ctx, "DELETE FROM users")
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
		elements := helpers.GetStructFields(models.User{}, []interface{}{"password", "otp_secret_key", "phone_no"})
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
		email = "invalid email"
		password = "invalid password"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"email", "password"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with an email that does not exist")
		email = "doesnotexist@test.com"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	It("should be an error", func() {
		By("sending a request with a password that does not match")
		password = "Notmatch@123"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
