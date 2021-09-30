package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

func CreateUserRow(ctx context.Context, option SQLOption) *SQLResponse {
	params := []string{}
	for i := 0; i < len(option.InsertColumns); i++ {
		params = append(params, fmt.Sprintf("$%v", i+1))
	}
	paramString := strings.Join(params, ", ")
	insertColumns := strings.Join(option.InsertColumns, ", ")
	returnColumns := strings.Join(option.ReturnColumns, ", ")

	sql := fmt.Sprintf(`INSERT INTO users (%v) VALUES (%v) RETURNING  %v`, insertColumns, paramString, returnColumns)
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, option.Arguments...).Scan(option.Destination...)
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

func FindUserRow(ctx context.Context, option SQLOption) *SQLResponse {
	returnColumns := strings.Join(option.ReturnColumns, ", ")
	sql := fmt.Sprintf(`SELECT %v FROM users %v`, returnColumns, option.AfterTableClauses)
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, option.Arguments...).Scan(option.Destination...)
	if errors.Is(err, pgx.ErrNoRows) {
		return &SQLResponse{
			StatusCode: http.StatusNotFound,
			Body:       gin.H{"message": "User not found"},
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
