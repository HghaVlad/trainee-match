package schemaregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

const (
	contentType = "application/vnd.schemaregistry.v1+json"
)

type Client struct {
	conf       *config.SchemaRegistry
	httpClient *http.Client
}

func NewClient(conf *config.SchemaRegistry) *Client {
	return &Client{conf: conf, httpClient: &http.Client{Timeout: conf.TimeOut}}
}

func (cl *Client) LookUpSchemaID(ctx context.Context, subject string, schema string) (int, error) {
	body, err := json.Marshal(schemaRequest{Schema: schema})
	if err != nil {
		return -1, err
	}
	url := fmt.Sprintf("%s/subjects/%s/versions", cl.conf.BaseURL, subject)

	httpRequeqst, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return -1, err
	}
	httpRequeqst.Header.Set("Content-Type", contentType)
	httpRequeqst.Header.Set("Accept", contentType)

	resp, err := cl.httpClient.Do(httpRequeqst)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return -1, fmt.Errorf("schema registry returned non-200 status code: %d", resp.StatusCode)
	}

	var response schemaVersionResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	return response.ID, nil

}

type schemaRequest struct {
	Schema string `json:"schema"`
}

type schemaVersionResponse struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Subject string `json:"subject"`
}
