package gql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Query performs a POST request with the given GraphQL query string
// and decodes the result into out.
func Query(ctx context.Context, url string, query string, out any) error {
	body, _ := json.Marshal(struct {
		Q string `json:"query"`
	}{Q: query})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("graphql status %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
