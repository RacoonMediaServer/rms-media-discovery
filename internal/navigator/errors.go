package navigator

import (
	"context"
	"errors"
	"path"

	uuid "github.com/satori/go.uuid"
)

func (p *Page) checkError(method string) {
	if p.err == nil {
		return
	}

	tmpUUID := uuid.NewV4().String()

	if p.batch != "" {
		p.log.Errorf("batch '%s' error: %s failed: %+v [ %s ]", p.batch, method, p.err, tmpUUID)
	} else {
		p.log.Errorf("browser error: %s failed: %+v [ %s ]", method, p.err, tmpUUID)
	}

	if errors.Is(p.err, context.DeadlineExceeded) || errors.Is(p.err, context.Canceled) {
		return
	}

	p.Screenshot(path.Join(p.dumpPath, tmpUUID+".jpg")).TracePage(path.Join(p.dumpPath, tmpUUID+".html"))
}
