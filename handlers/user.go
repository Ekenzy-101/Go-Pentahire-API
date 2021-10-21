package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/models"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func GetUser(c *gin.Context) {
	userId := c.Param("id")
	_, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User with the given id is invalid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `
	WITH cte_vehicles AS (
  	SELECT id, 
			average_rating, 
			image,
			make,
			name, 
			rental_fee,
			trips_count 
		FROM vehicles
		WHERE user_id = $1
	)
	SELECT u.id, 
		u.average_rating, 
		u.created_at,
		u.firstname,
		u.image,
		CASE 
			WHEN email_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
			ELSE CAST('true' AS BOOLEAN)
		END AS is_email_verified,
		CASE 
			WHEN phone_verified_at IS NULL THEN CAST ('false' AS BOOLEAN)
			ELSE CAST('true' AS BOOLEAN)
		END AS is_phone_verified,
		u.reviews_count,
		u.trips_count,
		json_agg(to_jsonb(v)) AS vehicles
	FROM users AS u, cte_vehicles AS v
	WHERE u.id = $1
	GROUP BY u.id
	`
	user := models.User{}
	destination := []interface{}{
		&user.ID,
		&user.AverageRating,
		&user.CreatedAt,
		&user.Firstname,
		&user.Image,
		&user.IsEmailVerified,
		&user.IsPhoneVerified,
		&user.ReviewsCount,
		&user.TripsCount,
		&user.Vehicles,
	}
	pool := services.GetPostgresConnectionPool()
	err = pool.QueryRow(ctx, sql, userId).Scan(destination...)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
