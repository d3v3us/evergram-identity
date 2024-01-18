package account

import (
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) GetByID(id string) (*Account, error) {
	var account Account
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetByEmail(email string) (*Account, error) {
	var account Account
	if err := r.db.Where("email = ?", email).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *AccountRepository) GetByUsername(name string) (*Account, error) {
	var account Account
	if err := r.db.Where("name = ?", name).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *AccountRepository) Create(account *Account) error {
	if err := r.db.Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Update(account *Account) error {
	if err := r.db.Save(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Delete(account *Account) error {
	if err := r.db.Delete(account).Error; err != nil {
		return err
	}
	return nil
}
