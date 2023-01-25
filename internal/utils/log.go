package utils

import (
	"context"
	"github.com/apex/log"
)

func LogFromContext(ctx context.Context, from string, defaultLog *log.Entry) *log.Entry {
	l, ok := ctx.Value("log").(*log.Entry)
	if ok {
		return l.WithField("from", from)
	}
	return defaultLog
}
