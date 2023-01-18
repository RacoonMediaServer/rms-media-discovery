package accounts

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"time"
)

const retryInterval = 24 * time.Hour

type account struct {
	model.Account
	marked     bool
	markedTime time.Time
}

type repository struct {
	mapIdToIndex map[string]int
	accounts     []account
	rrIndex      int
}

func newRepository() *repository {
	return &repository{
		mapIdToIndex: map[string]int{},
	}
}

func (r *repository) Add(acc model.Account) {
	r.accounts = append(r.accounts, account{Account: acc})
	r.mapIdToIndex[acc.Id] = len(r.accounts) - 1
}

func (r *repository) Delete(id string) error {
	idx, ok := r.mapIdToIndex[id]
	if !ok {
		return ErrNotFound
	}

	delete(r.mapIdToIndex, id)
	r.accounts = append(r.accounts[:idx], r.accounts[idx+1:]...)
	r.rrIndex = 0

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

	return r.accounts[idx].Account, true
}

func (r *repository) nextIndex() int {
	cur := r.rrIndex
	found := -1
	for found == -1 {
		acc := &r.accounts[r.rrIndex]
		if acc.marked {
			if time.Since(acc.markedTime) >= retryInterval {
				acc.marked = false
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
	r.accounts[idx].marked = true
	r.accounts[idx].markedTime = time.Now()
}
