package tests

import (
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

var _ = Describe("GET /auth/me", func() {
	var (
		accessToken  string
		userId       string
		responseBody gin.H
	)

	var ExecuteRequest = func() (*httptest.ResponseRecorder, error) {
		request, err := http.NewRequest(http.MethodGet, "/auth/me", nil)
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
		options := models.SQLOptions{
			Arguments:     []interface{}{"Test", "Test", "Test", "Test"},
			InsertColumns: []string{"email", "firstname", "lastname", "password"},
			Destination:   []interface{}{&userId},
			ReturnColumns: []string{"id"},
		}
		sqlResponse := models.InsertUserRow(ctx, options)
		Expect(sqlResponse).To(BeNil())

		user := &models.User{ID: userId}
		token, err := user.GenerateAccessToken()
		Expect(err).NotTo(HaveOccurred())

		accessToken = token
		responseBody = gin.H{}
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

		By("returning a body that contains the user's info")
		actual := helpers.GetMapKeys(responseBody["user"])
		elements := helpers.GetStructFields(models.User{}, []interface{}{"password", "otp_secret_key", "phone_no"})
		Expect(actual).To(ContainElements(elements...))
	})

	It("should be a success", func() {
		By("sending a request with an invalid or expired access token")
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

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains no user's info")
		Expect(responseBody["user"]).To(BeNil())
	})

	It("should be a success", func() {
		By("sending a request with an access token of a user that does not exist")
		_, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userId)
		Expect(err).NotTo(HaveOccurred())

		response, err := ExecuteRequest()
		Expect(err).NotTo(HaveOccurred())

		By("returning a status code of 200")
		Expect(response).To(HaveHTTPStatus(http.StatusOK))

		By("returning a body that contains no user's info")
		Expect(responseBody["user"]).To(BeNil())
	})
})
