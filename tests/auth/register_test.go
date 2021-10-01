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

var _ = Describe("POST /auth/register", func() {
	var (
		email        string
		firstname    string
		lastname     string
		password     string
		responseBody gin.H
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{
			"email":     email,
			"password":  password,
			"firstname": firstname,
			"lastname":  lastname,
		}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(requestBodyBytes))
		if err != nil {
			return nil, err
		}

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
		email = "test@test.com"
		firstname = "Test"
		lastname = "Test"
		password = "Testing@123"
		responseBody = gin.H{}
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

		By("returning a body that contains the user's info")
		actual := helpers.GetMapKeys(responseBody["user"])
		elements := helpers.GetStructFields(models.User{}, []interface{}{"password", "otp_secret_key"})
		Expect(actual).To(ContainElements(elements...))

		By("returning cookies")
		Expect(response.Result().Header).To(HaveKey("Set-Cookie"))
	})

	It("should be an error", func() {
		By("sending a request with invalid inputs")
		email = "invalid email"
		password = "invalid password"
		firstname = "1234"
		lastname = "1234"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"email", "firstname", "lastname", "password"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with an email that already exists")
		options := models.SQLOptions{
			Arguments:     []interface{}{email, firstname, lastname, password},
			InsertColumns: []string{"email", "firstname", "lastname", "password"},
			ReturnColumns: []string{"email"},
			Destination:   []interface{}{&email},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})
})
