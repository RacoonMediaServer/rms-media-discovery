package requester

import (
	"context"
	"encoding/json"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
	"time"
)

const Timeout = 30 * time.Second

var httpClient = http.Client{Timeout: Timeout}

type Requester interface {
	Get(ctx context.Context, url string, response interface{}) error
	Download(ctx context.Context, url string) ([]byte, error)
}

type requester struct {
	p provider.Provider
}

func New(p provider.Provider) Requester {
	return &requester{p: p}
}

func (r requester) Get(ctx context.Context, url string, response interface{}) error {
	buf, err := r.Download(ctx, url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, response); err != nil {
		return fmt.Errorf("'%s': unmarshal response failed: %w", url, err)
	}

	return nil
}

func (r requester) Download(ctx context.Context, url string) ([]byte, error) {
	l := utils.LogFromContext(ctx, r.p.ID(), log.WithField("from", r.p.ID())).WithField("url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req = req.WithContext(ctx)

	timer := prometheus.NewTimer(OutgoingRequestsMetric.WithLabelValues(r.p.ID()))

	var status int
	defer func() {
		OutgoingRequestsCounter.WithLabelValues(fmt.Sprintf("%d", status), r.p.ID()).Inc()
		timer.ObserveDuration()
	}()

	l.Debugf("Fetching...")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("'%s' request failed: %w", url, err)
	}
	defer resp.Body.Close()
	status = resp.StatusCode

	l.Debugf("Got response: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("'%s': unexpected status code: %d", url, resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("'%s': network I/O error: %w", url, err)
	}

	return buf, nil
}
