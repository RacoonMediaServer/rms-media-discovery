package accounts

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/apex/log"
	"time"
)

const retryInterval = 24 * time.Hour

type account struct {
	model.Account
	marked     bool
	markedTime time.Time
	reqPerDay  uint64
}

type repository struct {
	mapIdToIndex map[string]int
	accounts     []account
	rrIndex      int
	log          *log.Entry
}

func newRepository(log *log.Entry) *repository {
	return &repository{
		mapIdToIndex: map[string]int{},
		log:          log,
	}
}

func (r *repository) Add(acc model.Account) {
	r.accounts = append(r.accounts, account{Account: acc})
	r.mapIdToIndex[acc.Id] = len(r.accounts) - 1
	r.log.Debugf("Account '%s' registered", acc.Id)
}

func (r *repository) Delete(id string) error {
	idx, ok := r.mapIdToIndex[id]
	if !ok {
		return ErrNotFound
	}

	r.accounts = append(r.accounts[:idx], r.accounts[idx+1:]...)
	r.rrIndex = 0

	r.mapIdToIndex = make(map[string]int)
	for i := range r.accounts {
		r.mapIdToIndex[r.accounts[i].Id] = i
	}

	r.log.Debugf("Account '%s' deleted", id)

	return nil
}

func (r *repository) Get() (model.Account, bool) {
	if len(r.accounts) == 0 {
		return model.Account{}, false
	}

	idx := r.nextIndex()
	if idx < 0 {
		return model.Account{}, false
	}

	acc := &r.accounts[idx]
	acc.reqPerDay++
	if acc.Limit != 0 && acc.reqPerDay >= uint64(acc.Limit) {
		acc.Mark()
		r.log.Infof("Account '%s' is marked unaccessible, because limit reached (%d / %d)", acc.Id, acc.reqPerDay, acc.Limit)
	}

	r.log.Debugf("Account '%s' is used", acc.Id)

	return acc.Account, true
}

func (r *repository) nextIndex() int {
	cur := r.rrIndex
	found := -1
	for found == -1 {
		acc := &r.accounts[r.rrIndex]
		if acc.marked {
			if time.Since(acc.markedTime) >= retryInterval {
				acc.UnMark()
				r.log.Infof("Account '%s' is recovered", acc.Id)
				found = r.rrIndex
			}
		} else {
			found = r.rrIndex
		}
		r.rrIndex++
		if r.rrIndex >= len(r.accounts) {
			r.rrIndex = 0
		}
		if r.rrIndex == cur {
			return found
		}
	}

	return found
}

func (r *repository) MarkUnaccessible(id string) {
	idx, ok := r.mapIdToIndex[id]
	if !ok {
		return
	}
	acc := &r.accounts[idx]
	acc.Mark()
	r.log.Infof("Account '%s' is marked unaccessible", acc.Id)
}

func (a *account) Mark() {
	a.marked = true
	a.markedTime = time.Now()
	unaccessibleAccountsGauge.WithLabelValues(a.Service()).Inc()
}

func (a *account) UnMark() {
	a.marked = false
	a.reqPerDay = 0
	unaccessibleAccountsGauge.WithLabelValues(a.Service()).Dec()
}
