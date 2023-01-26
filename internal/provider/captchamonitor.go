package provider

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	captchaSolvedCounter    *prometheus.CounterVec
	captchaNotSolvedCounter *prometheus.CounterVec
)

func init() {
	captchaSolvedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "captcha_solved_count",
		Help:      "Amount of solved captchas",
	}, []string{"provider"})

	captchaNotSolvedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "captcha_not_solved_count",
		Help:      "Amount of not solved captchas",
	}, []string{"provider"})
}

type captchaMonitor struct {
	solver CaptchaSolver
}

func (c captchaMonitor) ID() string {
	return c.solver.ID()
}

func (c captchaMonitor) Solve(ctx context.Context, captcha Captcha) (string, error) {
	result, err := c.solver.Solve(ctx, captcha)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			captchaNotSolvedCounter.WithLabelValues(c.solver.ID()).Inc()
		}
	} else {
		captchaSolvedCounter.WithLabelValues(c.solver.ID()).Inc()
	}

	return result, err
}

func NewCaptchaSolverMonitor(solver CaptchaSolver) CaptchaSolver {
	return &captchaMonitor{solver: solver}
}
