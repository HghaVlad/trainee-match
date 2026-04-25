package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

func GetAuthClient(ctx context.Context, authServiceBaseURL string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	logger := slog.Default()

	username := "usertest" + strings.ReplaceAll(uuid.NewString(), "-", "")
	email := username + "@gmail.com"

	registerBody := fmt.Sprintf(`{
		"first_name":"Test",
		"last_name":"User",
		"email":"%s",
		"username":"%s",
		"password":"testpass",
		"role":"%s"
	}`, email, username, identity.RoleHR)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authServiceBaseURL+"/auth/register",
		strings.NewReader(registerBody),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create register user request", "err", err)
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		logger.ErrorContext(ctx, "failed to register user", "err", err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		logger.ErrorContext(ctx,
			"negative register response",
			"code", resp.StatusCode,
			"status", resp.Status,
			"err_msg", errResp.Error,
		)
		return nil, fmt.Errorf("register user status 200 expected, returned %d", resp.StatusCode)
	}

	loginBody := fmt.Sprintf(`{
		"username":"%s",
		"password":"testpass"
	}`, username)

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authServiceBaseURL+"/auth/login",
		strings.NewReader(loginBody),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create login user request", "err", err)
		return nil, err
	}

	resp, err = client.Do(req)

	if err != nil {
		logger.ErrorContext(ctx, "failed to login user", "err", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorContext(ctx,
			"negative login response",
			"code", resp.StatusCode,
			"status", resp.Status,
		)
		return nil, err
	}

	for _, c := range resp.Cookies() {
		logger.InfoContext(ctx, "cookie:", "name", c.Name, "val", c.Value)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return client, nil
}
