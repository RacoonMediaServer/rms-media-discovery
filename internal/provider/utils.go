package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io"
	"net/http"
)

func doRequest(l *log.Entry, cli http.Client, ctx context.Context, url string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req = req.WithContext(ctx)

	l.Debugf("Fetching '%s'...", url)
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	l.Debugf("'%s' response: %s", url, resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("network I/O error: %w", err)
	}

	if err := json.Unmarshal(buf, response); err != nil {
		return fmt.Errorf("unmarshal response from failed: %w", err)
	}

	return nil
}
