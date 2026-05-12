package schemaregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ContentType    = "application/vnd.schemaregistry.v1+json"
	DefaultTimeout = 10 * time.Second
)

type schemaRequest struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType,omitempty"`
}

type schemaVersionResponse struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Subject string `json:"subject"`
}

type RealRegistryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *RealRegistryClient {
	return &RealRegistryClient{baseURL: baseURL, httpClient: &http.Client{Timeout: DefaultTimeout}}
}

func (client *RealRegistryClient) RegisterSchema(ctx context.Context, subject, schema string) (int, error) {
	url := fmt.Sprintf("%s/subjects/%s/versions", client.baseURL, subject)

	request := schemaRequest{
		Schema: schema,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", ContentType)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, errors.New("schema registration failed: " + resp.Status)
	}
	var response schemaVersionResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	if response.ID == 0 {
		return 0, errors.New("schema registration failed")
	}
	return response.ID, nil
}
