package tests

import (
	"testing"

	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuthRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth")
}

var (
	pool *pgxpool.Pool

	_ = BeforeSuite(func() {
		pool = services.CreatePostgresConnectionPool()
	})

	_ = AfterSuite(func() {
		pool.Close()
	})
)
