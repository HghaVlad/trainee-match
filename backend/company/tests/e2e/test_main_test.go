package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

	// Postgres

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

	if err != nil {
		log.Fatalf("Error creating postgres container: %v", err)
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Error getting postgres host: %v", err)
	}
	pgPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("Error getting postgres port: %v", err)
	}
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

	// Redis

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

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Printf("failed to terminate redis container: %s", err)
		}
	}()

	redisHost, err := redisC.Host(ctx)
	if err != nil {
		log.Fatalf("Error getting redis host: %v", err)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		log.Fatalf("Error getting redis port: %v", err)
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())

	log.Println("redis:", redisAddr)

	// Keycloak

	keycloakRealmImport, err := filepath.Abs("../../../auth/import/trainee-match-realm.json")
	if err != nil {
		log.Fatalf("resolve keycloak realm import path: %v", err)
	}

	keycloakContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "quay.io/keycloak/keycloak:26.5",
			Env:          map[string]string{"KEYCLOAK_ADMIN": "admin", "KEYCLOAK_ADMIN_PASSWORD": "admin"},
			ExposedPorts: []string{"8080/tcp"},
			Cmd:          []string{"start-dev", "--import-realm"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      keycloakRealmImport,
					ContainerFilePath: "/opt/keycloak/data/import/trainee-match-realm.json",
					FileMode:          0o644,
				},
			},
			WaitingFor: wait.ForAll(
				wait.ForHTTP("/realms/trainee-match").WithPort("8080/tcp").WithStartupTimeout(6*time.Minute),
				wait.ForLog("Running the server").WithStartupTimeout(6*time.Minute),
			),
		},
		Started: true,
		Logger:  log.New(os.Stdout, "", log.LstdFlags),
	})
	if err != nil {
		log.Fatalf("Error creating keycloak container: %v", err)
	}

	defer func() {
		if err := keycloakContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate keycloak container: %s", err)
		}
	}()

	keycloakHost, err := keycloakContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Error getting keycloak host: %v", err)
	}
	keycloakPort, err := keycloakContainer.MappedPort(ctx, "8080/tcp")
	if err != nil {
		log.Fatalf("Error getting keycloak port: %v", err)
	}
	keycloakExternalURL := "http://" + net.JoinHostPort(keycloakHost, keycloakPort.Port())

	// Config

	jwkURL := strings.TrimRight(keycloakExternalURL, "/") + "/realms/trainee-match/protocol/openid-connect/certs"

	conf := &config.Config{
		HTTP: config.HTTPConfig{
			JWKUrl: jwkURL,
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
		log.Fatal("couldn't build app:", err)
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
