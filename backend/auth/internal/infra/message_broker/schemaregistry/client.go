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
	contentType    = "application/vnd.schemaregistry.v1+json"
	defaultTimeout = 10 * time.Second
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
	baseUrl    string
	httpClient *http.Client
}

func NewClient(baseUrl string) *RealRegistryClient {
	return &RealRegistryClient{baseUrl: baseUrl, httpClient: &http.Client{Timeout: defaultTimeout}}
}

func (client *RealRegistryClient) RegisterSchema(ctx context.Context, subject, schema string) (int, error) {
	url := fmt.Sprintf("%s/subjects/%s/versions", client.baseUrl, subject)

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
	req.Header.Set("Content-Type", contentType)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, errors.New("Schema registration failed: " + resp.Status)
	}
	var response schemaVersionResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	if response.ID == 0 {
		return 0, errors.New("Schema registration failed")
	}
	return response.ID, nil
}
