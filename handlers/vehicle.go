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

func GetVehicle(c *gin.Context) {
	vehicleId := c.Param("id")
	_, err := uuid.Parse(vehicleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vehicle with the given id is invalid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `
	WITH cte_users AS (
  	SELECT id, firstname, created_at, image, trips_count FROM users
	)
	SELECT v.id, 
		v.average_rating, 
		v.created_at,
		v.is_rented,
		v.image,
		v.location,
		v.make,
		v.name,
		v.rental_fee,
		v.reviews_count,
		to_jsonb(u) as user, 
		v.trips_count
	FROM vehicles AS v 
	JOIN cte_users AS u ON v.user_id = u.id 
	WHERE v.id = $1`
	vehicle := &models.Vehicle{Location: &models.Location{}}
	destination := []interface{}{
		&vehicle.ID,
		&vehicle.AverageRating,
		&vehicle.CreatedAt,
		&vehicle.IsRented,
		&vehicle.Image,
		vehicle.Location,
		&vehicle.Make,
		&vehicle.Name,
		&vehicle.RentalFee,
		&vehicle.ReviewsCount,
		&vehicle.User,
		&vehicle.TripsCount,
	}
	pool := services.GetPostgresConnectionPool()
	err = pool.QueryRow(ctx, sql, vehicleId).Scan(destination...)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vehicle not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicle": vehicle})
}
