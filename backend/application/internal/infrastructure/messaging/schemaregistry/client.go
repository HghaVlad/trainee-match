package schemaregistry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

const (
	retryCount  = 3
	retryDelay  = 200 * time.Millisecond
	maxWaitTime = 2 * time.Second
	contentType = "application/vnd.schemaregistry.v1+json"
)

var (
	ErrSchemaNotFound            = errors.New("schema not found")
	ErrSchemaRegistryUnavailable = errors.New("schema registry is unavailable now")
)

type Client struct {
	conf  config.SchemaRegistry
	resty *resty.Client
}

func NewClient(conf config.SchemaRegistry) *Client {
	return &Client{
		conf: conf,
		resty: resty.New().
			SetTimeout(conf.TimeOut).
			SetRetryCount(retryCount).
			SetRetryWaitTime(retryDelay).
			SetRetryMaxWaitTime(maxWaitTime),
	}
}

func (cl *Client) LookUpSchemaID(
	ctx context.Context,
	subject string,
	schema string,
) (int, error) {
	var success schemaVersionResponse
	var apiErr apiError

	reqBody := schemaRequest{
		Schema: schema,
	}

	resp, err := cl.resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", contentType).
		SetHeader("Accept", contentType).
		SetBody(reqBody).
		SetResult(&success).
		SetError(&apiErr).
		Post(fmt.Sprintf("%s/subjects/%s", cl.conf.BaseURL, subject))

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

func (cl *Client) GetSchemaByID(ctx context.Context, id int) (string, error) {
	var success schemaByIDResponse
	var apiErr apiError

	resp, err := cl.resty.R().
		SetContext(ctx).
		SetHeader("Accept", contentType).
		SetResult(&success).
		SetError(&apiErr).
		Get(fmt.Sprintf("%s/schemas/ids/%d", cl.conf.BaseURL, id))

	if err != nil {
		return "", fmt.Errorf("schema reg client: get schema by id: %w", err)
	}

	if resp.IsError() {
		err := handleRestyError(resp, &apiErr)
		return "", err
	}

	return success.Schema, nil
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
	Schema string `json:"schema"`
}

type schemaVersionResponse struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Subject string `json:"subject"`
}

type schemaByIDResponse struct {
	Schema string `json:"schema"`
}

type apiError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
