package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	_ "path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	appl "github.com/HghaVlad/trainee-match/backend/company/internal/app"
	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/tests/e2e/helpers"
)

var (
	AuthClient *http.Client // logged in client
	app        *appl.App
	baseURL    string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "test_db"
	dbUser := "test_user"
	dbPass := "test_pass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.BasicWaitStrategies(),
	)

	pgHost, _ := postgresContainer.Host(ctx)
	pgPort, _ := postgresContainer.MappedPort(ctx, "5432/tcp")
	pgSPort := pgPort.Port()

	log.Printf("postgres host: %s, port: %s", pgHost, pgSPort)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		pgHost,
		pgSPort,
		dbName,
	)

	migErr := runMigrations(dbURL)
	if migErr != nil {
		log.Fatalln(errors.Unwrap(migErr))
	}

	redisC, err := testcontainers.Run(
		ctx, "redis:latest",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	redisHost, _ := redisC.Host(ctx)
	redisPort, _ := redisC.MappedPort(ctx, "6379/tcp")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())

	log.Println(redisAddr)

	conf := &config.Config{
		HTTP: config.HTTPConfig{
			JWKUrl: "http://localhost:8080/realms/trainee-match/protocol/openid-connect/certs",
		},
		CompanyDB: config.DBConfig{
			Host:            pgHost,
			Port:            pgSPort,
			Name:            dbName,
			User:            dbUser,
			Password:        dbPass,
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 30 * time.Second,
		},
		Redis: config.RedisConfig{
			Host: redisHost,
			Port: redisPort.Port(),
		},
	}

	app, err = appl.Build(conf)
	if err != nil {
		log.Fatal(err)
	}

	AuthClient = helpers.GetAuthClient()

	server := httptest.NewServer(app.HttpSrv.Handler)
	baseURL = server.URL

	code := m.Run()

	server.Close()
	os.Exit(code)
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://../../internal/infrastructure/db/postgres/migrations",
		dbURL,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
