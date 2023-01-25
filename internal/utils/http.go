package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io"
	"net/http"
	"time"
)

const requestTimeout = 2 * time.Minute

var HttpClient = http.Client{Timeout: requestTimeout}

func Get(l *log.Entry, ctx context.Context, url string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req = req.WithContext(ctx)

	l.Debugf("Fetching '%s'...", url)
	resp, err := HttpClient.Do(req)
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

func Download(l *log.Entry, ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req = req.WithContext(ctx)

	l.Debugf("Downloading '%s'...", url)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	l.Debugf("'%s' response: %s", url, resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("network I/O error: %w", err)
	}

	return buf, nil
}
