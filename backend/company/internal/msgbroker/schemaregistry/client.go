package schemaregistry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/go-resty/resty/v2"
)

const (
	defaultTimeout = 10 * time.Second
	retryCount     = 3
	retryDelay     = 200 * time.Millisecond
	maxWaitTime    = 2 * time.Second
	contentType    = "application/vnd.schemaregistry.v1+json"
)

var (
	ErrSchemaNotFound            = errors.New("schema not found")
	ErrSchemaRegistryUnavailable = errors.New("schema registry is unavailable now")
)

type RealRegistryClient struct {
	resty   *resty.Client
	baseURL string
}

func NewClient(cfg config.SchemaRegistry) *RealRegistryClient {
	return &RealRegistryClient{
		resty: resty.New().
			SetTimeout(defaultTimeout).
			SetRetryCount(retryCount).
			SetRetryWaitTime(retryDelay).
			SetRetryMaxWaitTime(maxWaitTime),
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
	}
}

func (c *RealRegistryClient) LookupSchemaID(
	ctx context.Context,
	subject string,
	schema string,
) (int, error) {
	var success schemaVersionResponse
	var apiErr apiError

	reqBody := schemaRequest{
		Schema: schema,
	}

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", contentType).
		SetHeader("Accept", contentType).
		SetBody(reqBody).
		SetResult(&success).
		SetError(&apiErr).
		Post(fmt.Sprintf("%s/subjects/%s", c.baseURL, subject))

	if err != nil {
		return 0, fmt.Errorf("schema reg client: lookup schema: %w: %w", ErrSchemaRegistryUnavailable, err)
	}

	if resp.IsError() {
		// 40401 → schema not found under subject
		if resp.StatusCode() == http.StatusNotFound && apiErr.ErrorCode == 40401 {
			return 0, ErrSchemaNotFound
		}

		return 0, handleRestyError(resp, &apiErr)
	}

	return success.ID, nil
}

func handleRestyError(resp *resty.Response, apiErr *apiError) error {
	if resp.StatusCode() == http.StatusNotFound && apiErr.ErrorCode == 40403 {
		return ErrSchemaNotFound
	}

	if apiErr.Message != "" {
		return fmt.Errorf("schema registry returned %d (%d): %s",
			resp.StatusCode(), apiErr.ErrorCode, apiErr.Message)
	}
	return fmt.Errorf("schema registry returned %d: %s",
		resp.StatusCode(), resp.String())
}

type schemaRequest struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType,omitempty"`
}

type schemaVersionResponse struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Subject string `json:"subject"`
}

type apiError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
