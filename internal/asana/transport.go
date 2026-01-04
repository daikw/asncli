package asana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) do(ctx context.Context, method, path string, query map[string]string, body any, out any) error {
	var env responseEnvelope
	if err := c.doRaw(ctx, method, path, query, body, &env); err != nil {
		return err
	}
	return json.Unmarshal(env.Data, out)
}

func (c *Client) doRaw(ctx context.Context, method, path string, query map[string]string, body any, out any) error {
	token, err := c.tokenSource.Token(ctx)
	if err != nil {
		return err
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}
	q := u.Query()
	for key, value := range query {
		if value != "" {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode >= 300 {
		var apiErr errorResponse
		bodyBytes, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return fmt.Errorf("failed to read error response: %w", errRead)
		}
		if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && len(apiErr.Errors) > 0 {
			return fmt.Errorf("response: %s", apiErr.Errors[0].Message)
		}
		return fmt.Errorf("unexpected response: %s", strings.TrimSpace(string(bodyBytes)))
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}

	return nil
}
