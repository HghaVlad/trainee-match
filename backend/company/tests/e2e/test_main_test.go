package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	appl "github.com/HghaVlad/trainee-match/backend/company/internal/app"
	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/tests/e2e/helpers"
)

//nolint:gochecknoglobals // tests don't have params
var (
	AuthClient         *http.Client // logged in client
	app                *appl.App
	baseURL            string
	authServiceBaseURL string
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

//nolint:gocognit // setups up docker containers, uses them for tests, and terminates them
func run(m *testing.M) int {
	ctx := context.Background()

	logger := slog.Default()

	// --------
	// Network
	// --------

	netwrk, err := network.New(ctx)
	if err != nil {
		logger.Error("failed to create docker network", "err", err)
		return 1
	}

	defer func() {
		if err := netwrk.Remove(ctx); err != nil {
			logger.Error("failed to remove docker network", "err", err)
		}
	}()

	// --------
	// Postgres
	// --------

	dbName := "test_db"
	dbUser := "test_user"
	dbPass := "test_pass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{"company-postgres"}, netwrk),
	)
	if err != nil {
		logger.Error("failed to start company postgres container", "err", err)
		return 1
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			logger.Warn("failed to terminate company postgres container", "err", err)
		}
	}()

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		logger.Error("failed to get company postgres host", "err", err)
		return 1
	}

	pgPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		logger.Error("failed to get company postgres port", "err", err)
		return 1
	}

	pgSPort := pgPort.Port()

	logger.Info("company postgres ready", "host", pgHost, "port", pgSPort)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPass, net.JoinHostPort(pgHost, pgSPort), dbName,
	)

	if err := runMigrations(dbURL); err != nil {
		logger.Error("failed to run company postgres migrations", "err", err)
		return 1
	}

	// -----
	// Redis
	// -----

	redisC, err := testcontainers.Run(
		ctx, "redis:latest",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		logger.Error("failed to start redis container", "err", err)
		return 1
	}

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			logger.Warn("failed to terminate redis container", "err", err)
		}
	}()

	redisHost, err := redisC.Host(ctx)
	if err != nil {
		logger.Error("failed to get redis host", "err", err)
		return 1
	}

	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		logger.Error("failed to get redis port", "err", err)
		return 1
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())

	logger.Info("redis ready", "addr", redisAddr)

	// ---------------------------
	// Keycloak
	// ---------------------------

	keycloakRealmImport, err := filepath.Abs("../../../auth/import/trainee-match-realm.json")
	if err != nil {
		logger.Error("failed to resolve keycloak realm import path", "err", err)
		return 1
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
			Networks: []string{netwrk.Name},
			NetworkAliases: map[string][]string{
				netwrk.Name: {"keycloak"},
			},
		},
		Started: true,
	})
	if err != nil {
		logger.Error("failed to start keycloak container", "err", err)
		return 1
	}

	defer func() {
		if err := keycloakContainer.Terminate(ctx); err != nil {
			logger.Warn("failed to terminate keycloak container", "err", err)
		}
	}()

	keycloakHost, err := keycloakContainer.Host(ctx)
	if err != nil {
		logger.Error("failed to get keycloak host", "err", err)
		return 1
	}

	keycloakPort, err := keycloakContainer.MappedPort(ctx, "8080/tcp")
	if err != nil {
		logger.Error("failed to get keycloak port", "err", err)
		return 1
	}

	keycloakExternalURL := "http://" + net.JoinHostPort(keycloakHost, keycloakPort.Port())
	keycloakInternalURL := "http://keycloak:8080"

	logger.Info("keycloak ready",
		"external_url", keycloakExternalURL,
		"internal_url", keycloakInternalURL,
	)

	// ------------
	// Auth Service
	// ------------

	authContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "ghcr.io/m0s1ck/backend-auth-service:1.0",
			Env:          map[string]string{"KC_URL": keycloakInternalURL},
			ExposedPorts: []string{"8000/tcp"},
			Networks:     []string{netwrk.Name},
			WaitingFor:   wait.ForLog("Service is starting").WithStartupTimeout(1 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		logger.Error("failed to start auth service", "err", err)
		return 1
	}

	defer func() {
		if err := authContainer.Terminate(ctx); err != nil {
			logger.Warn("failed to terminate auth container", "err", err)
		}
	}()

	time.Sleep(30 * time.Second)

	authHost, err := authContainer.Host(ctx)
	if err != nil {
		logger.Error("failed to get auth host", "err", err)
		return 1
	}

	authPort, err := authContainer.MappedPort(ctx, "8000/tcp")
	if err != nil {
		logger.Error("failed to get auth port", "err", err)
		return 1
	}

	logger.Info("auth service ready", "host", authHost, "port", authPort.Port())

	// --------------------
	// Build App, Run Tests
	// --------------------

	authServiceBaseURL = fmt.Sprintf("http://%s/api/v1", net.JoinHostPort(authHost, authPort.Port()))

	jwkURL := strings.TrimRight(keycloakExternalURL, "/") + "/realms/trainee-match/protocol/openid-connect/certs"

	conf := &config.Config{
		HTTP: config.HTTP{
			JWKUrl: jwkURL,
		},
		Postgres: config.Postgres{
			Host:         pgHost,
			Port:         pgSPort,
			Name:         dbName,
			User:         dbUser,
			Password:     dbPass,
			SSLMode:      "disable",
			MaxPoolConns: 10,
			MinPoolConns: 2,
		},
		Redis: config.Redis{
			Host: redisHost,
			Port: redisPort.Port(),
		},
	}

	app, err = appl.Build(ctx, conf, logger)
	if err != nil {
		logger.Error("failed to build app", "err", err)
		return 1
	}

	AuthClient, err = helpers.GetAuthClient(ctx, authServiceBaseURL)
	if err != nil {
		return 1
	}

	server := httptest.NewServer(app.HTTPSrv.Handler)
	baseURL = server.URL

	logger.Info("test environment ready")

	code := m.Run()

	server.Close()
	return code
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
