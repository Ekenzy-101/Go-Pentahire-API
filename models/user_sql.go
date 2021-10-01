package models

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

func InsertUserRow(ctx context.Context, options SQLOptions) *SQLResponse {
	options.TableName = config.UsersTable
	options.Statement = InsertStatement

	sql := buildQuery(options)
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, options.Arguments...).Scan(options.Destination...)
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

func SelectUserRow(ctx context.Context, options SQLOptions) *SQLResponse {
	options.Statement = SelectStatement
	options.TableName = config.UsersTable

	sql := buildQuery(options)
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, options.Arguments...).Scan(options.Destination...)
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

func UpdateAndReturnUserRow(ctx context.Context, options SQLOptions) *SQLResponse {
	options.Statement = UpdateStatement
	options.TableName = config.UsersTable

	sql := buildQuery(options)
	pool := services.GetPostgresConnectionPool()
	err := pool.QueryRow(ctx, sql, options.Arguments...).Scan(options.Destination...)
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
