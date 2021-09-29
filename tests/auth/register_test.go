package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/routes"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /auth/register", func() {
	var (
		email        string
		firstname    string
		lastname     string
		password     string
		pool         *pgxpool.Pool
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

	BeforeSuite(func() {
		pool = services.CreatePostgresConnectionPool()
	})

	BeforeEach(func() {
		email = "test@test.com"
		firstname = "Test"
		lastname = "Test"
		password = "Testing@123"
		responseBody = gin.H{}
	})

	AfterEach(func() {
		_, err := pool.Exec(context.Background(), "DELETE FROM users; DELETE FROM verify_email;")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be a success", func() {
		By("sending a request with valid inputs")
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body which contains the user's info")
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

		By("returning a body which contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"email", "firstname", "lastname", "password"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with an email that already exists")
		sqlResponse := models.CreateUserRow(context.Background(), &models.User{
			Email:     email,
			Password:  password,
			Firstname: firstname,
			Lastname:  lastname,
		})
		Expect(sqlResponse).To(BeNil())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body which contains error messages")
		Expect(responseBody).To(HaveKey("message"))
	})

	AfterSuite(func() {
		pool.Close()
	})
})

func TestRegisterSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Register Suite")
}
