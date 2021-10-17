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

var _ = Describe("PUT /account/profile", func() {
	var (
		accessToken  string
		email        string
		firstname    string
		lastname     string
		responseBody gin.H
		userId       string
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		requestBodyMap := gin.H{"email": email, "firstname": firstname, "lastname": lastname}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPut, "/account/profile", bytes.NewReader(requestBodyBytes))
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
		firstname = "Test"
		lastname = "Test"
		email = "newemail@test.com"
		responseBody = gin.H{}
	})

	JustBeforeEach(func() {
		user := &models.User{}
		err := user.HashPassword()
		Expect(err).NotTo(HaveOccurred())

		options := models.SQLOptions{
			Arguments:     []interface{}{"verified@test.com", "Test", "Test", "Test", time.Now()},
			InsertColumns: []string{"email", "password", "firstname", "lastname", "email_verified_at"},
			ReturnColumns: []string{"id"},
			Destination:   []interface{}{&userId},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		user.ID = userId
		accessToken, err = user.GenerateAccessToken()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_, err := pool.Exec(ctx, "DELETE FROM users")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be a success", func() {
		By("sending a request with a valid access token and inputs")
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains a success message")
		Expect(responseBody).To(HaveKey("message"))

		By("updating user profile accordingly")
		var newFirstName, newLastName, newEmail string
		var emailVerifiedAt interface{}
		destination := []interface{}{&newFirstName, &newLastName, &newEmail, &emailVerifiedAt}
		sql := "SELECT firstname, lastname, email, email_verified_at FROM users WHERE id = $1"
		err = pool.QueryRow(ctx, sql, userId).Scan(destination...)
		Expect(err).NotTo(HaveOccurred())

		Expect(firstname).To(Equal(newFirstName))
		Expect(lastname).To(Equal(newLastName))
		Expect(email).To(Equal(newEmail))
		Expect(emailVerifiedAt).To(BeNil())
	})

	It("should be an error", func() {
		By("sending a request with invalid inputs")
		lastname = ""
		firstname = ""
		email = "invalid email"
		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 400")
		Expect(response).To(HaveHTTPStatus(http.StatusBadRequest))

		By("returning a body that contains error messages")
		actual := helpers.GetMapKeys(responseBody)
		elements := []interface{}{"email", "firstname", "lastname"}
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be an error", func() {
		By("sending a request with an email address that does exist")
		email = "doesexist@test.com"
		_, err := pool.Exec(ctx, "INSERT INTO users (email, firstname, lastname, password) VALUES ($1, '', '', '')", email)
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
