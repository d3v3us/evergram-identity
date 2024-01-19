package account

import (
	"github.com/deveusss/evergram-core/database"
)

type AccountRepository struct {
	db *database.OrmDatabase
}

func NewAccountRepository(db *database.OrmDatabase) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) GetByID(id string) (*Account, error) {
	var account Account
	if err := r.db.Orm.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetByEmail(email string) (*Account, error) {
	var account Account
	if err := r.db.Orm.Where("email = ?", email).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *AccountRepository) GetByUsername(name string) (*Account, error) {
	var account Account
	if err := r.db.Orm.Where("name = ?", name).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *AccountRepository) Create(account *Account) error {
	if err := r.db.Orm.Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Update(account *Account) error {
	if err := r.db.Orm.Save(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Delete(account *Account) error {
	if err := r.db.Orm.Delete(account).Error; err != nil {
		return err
	}
	return nil
}
