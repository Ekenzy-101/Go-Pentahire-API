package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Obj should be a pointer to a value
func ValidateRequestBody(c *gin.Context, obj interface{}) interface{} {
	err := c.ShouldBindJSON(obj)
	validationErrors := validator.ValidationErrors{}
	if errors.As(err, &validationErrors) {
		return GenerateErrorMessages(validationErrors)
	}

	if err != nil {
		return gin.H{"message": err.Error()}
	}

	return nil
}
