package accounts

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type account struct {
	model.Account
}

type repository struct {
	accounts map[string]account
}

func newRepository() *repository {
	return &repository{
		accounts: map[string]account{},
	}
}

func (r *repository) Add(acc model.Account) {
	r.accounts[acc.Id] = account{Account: acc}
}

func (r *repository) Delete(id string) error {
	_, ok := r.accounts[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.accounts, id)
	return nil
}
