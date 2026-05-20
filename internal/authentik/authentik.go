package authentik

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// An Authentik API client
type Authentik struct {
	client *http.Client
	host   string
	token  string
}

type ApiError struct {
	StatusCode int
	Body       map[string]any
}

func New(host, token string, caCert []byte, insecureSkipVerify bool) (*Authentik, error) {
	client, err := buildHttpClient(caCert, insecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	return &Authentik{
		client,
		host,
		token,
	}, nil
}

func buildHttpClient(caCert []byte, insecureSkipVerify bool) (*http.Client, error) {
	cfg := &tls.Config{InsecureSkipVerify: insecureSkipVerify}

	if len(caCert) > 0 {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		cfg.RootCAs = pool
	}

	return &http.Client{Transport: &http.Transport{TLSClientConfig: cfg}}, nil
}

func (c *Authentik) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.host+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed when executing request: %w", err)
	}

	return resp, nil
}

func (e ApiError) Error() string {
	body, err := json.Marshal(e.Body)
	if err != nil {
		return fmt.Sprintf("authentik API error: status %d", e.StatusCode)
	}
	return fmt.Sprintf("authentik API error: status %d, body: %s", e.StatusCode, body)
}

func parseApiError(resp *http.Response) error {
	var body map[string]any
	// Ignore decode error as having only the status code is still useful
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return &ApiError{StatusCode: resp.StatusCode, Body: body}
}
