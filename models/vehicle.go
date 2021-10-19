package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	ID            string    `json:"id"`
	Address       string    `json:"address,omitempty"`
	AverageRating float64   `json:"average_rating"`
	CreatedAt     time.Time `json:"created_at"`
	Image         string    `json:"image"`
	IsRented      bool      `json:"is_rented"`
	Make          string    `json:"make"  binding:"required"`
	Name          string    `json:"name" binding:"required"`
	Location      *Location `json:"location,omitempty"`
	RentalFee     int       `json:"rental_fee"`
	ReviewsCount  int       `json:"reviews_count"`
	TripsCount    int       `json:"trips_count"`
	User          gin.H     `json:"user,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
}
