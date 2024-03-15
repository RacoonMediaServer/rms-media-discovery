package db

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func (d *Database) LoadAccounts() (result []model.Account, err error) {
	result = make([]model.Account, 0)
	if err = d.conn.Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (d *Database) CreateAccount(acc model.Account) error {
	return d.conn.Create(acc).Error
}

func (d *Database) DeleteAccount(id string) error {
	return d.conn.Model(&model.Account{}).Unscoped().Delete(&model.Account{Id: id}).Error
}
