package tests

import (
	"context"
	"testing"

	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVerificationRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Verification")
}

var (
	pool        *pgxpool.Pool
	redisClient *redis.Client
	ctx         = context.Background()

	_ = BeforeSuite(func() {
		pool = services.CreatePostgresConnectionPool()
		redisClient = services.CreateRedisClient()
	})

	_ = AfterSuite(func() {
		pool.Close()
	})
)
