package models

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func CreateUserRow(ctx context.Context, user *User) *SQLResponse {
	args := []interface{}{strings.ToLower(user.Email), user.Firstname, user.Lastname, user.Password}
	sql := `INSERT INTO users (email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING id`
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, args...).Scan(&user.ID)
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return &SQLResponse{
			StatusCode: http.StatusBadRequest,
			Body:       gin.H{"message": "A user with the given email already exists"},
		}
	}

	if err != nil {
		return &SQLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       gin.H{"message": err.Error()},
		}
	}

	return nil
}

func CreateVerifyEmailRow(ctx context.Context, args []interface{}) *SQLResponse {
	sql := `INSERT INTO verify_email (user_id, token) VALUES ($1, $2)`
	pool := services.GetPostgresConnectionPool()
	_, err := pool.Exec(ctx, sql, args...)
	if err != nil {
		return &SQLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       gin.H{"message": err.Error()},
		}
	}

	return nil
}
