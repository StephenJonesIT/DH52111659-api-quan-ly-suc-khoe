package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) (error)
	Update(ctx context.Context, cond map[string]interface{}, updateValue map[string]interface{}) (error)
	GetByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountById(ctx context.Context, id string) (*models.Account, error)
	GetDB() *gorm.DB
	GetAccounts(ctx context.Context, paging *common.Paging,cond map[string]interface{}) ([]*models.Account, error)
	DeactivateAccount(ctx context.Context, tx *gorm.DB,id string) error
}


// AccountRepoImpl implements the AccountRepository interface
// using GORM for PostgreSQL
// and provides methods for account management.
// It includes methods for creating, updating, and retrieving accounts.
type AccountRepoImpl struct {
	DB *gorm.DB // Assuming you have a DB type that wraps your database connection
}

func NewAccountRepoImpl(db *gorm.DB) *AccountRepoImpl {
	return &AccountRepoImpl{DB: db}
}

func(repo *AccountRepoImpl) Create(ctx context.Context, account *models.Account) (error) {
	if err := repo.DB.WithContext(ctx).
		Table(models.Account{}.TableName()).
		Create(account).Error; err != nil {
		return err 
	}

	return nil
}

func(repo *AccountRepoImpl) Update(ctx context.Context, cond map[string]interface{}, updateValue map[string]interface{}) (error){
	if err := repo.DB.WithContext(ctx).
		Table(models.Account{}.TableName()).Where(cond).Updates(updateValue).Error; err != nil {
		return err 
	}
	return nil
}

func(repo *AccountRepoImpl) GetByEmail(ctx context.Context, email string) (*models.Account, error){
	var account models.Account

	if err := repo.DB.WithContext(ctx).
		Table(models.Account{}.TableName()).
		Where("email = ?", email).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		return nil, err // Other error
	}

	return &account, nil
}

func(repo *AccountRepoImpl) GetAccountById(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	if err := repo.DB.WithContext(ctx).
		Table(models.Account{}.TableName()).
		Where("id = ?", id).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil , nil
		}
		return nil , err 
	}

	return &account, nil
}

func (repo *AccountRepoImpl) GetDB() *gorm.DB{
	return repo.DB
}

func(repo *AccountRepoImpl) GetAccounts(ctx context.Context, paging *common.Paging,cond map[string]interface{}) ([]*models.Account, error) {
	var accounts []*models.Account

	query := repo.DB.WithContext(ctx).Table(models.Account{}.TableName()).Where(cond)
	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Order("created_at DESC").
		Offset((paging.Page-1)*paging.Limit).
		Limit(paging.Limit).
		Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func(repo *AccountRepoImpl) DeactivateAccount(ctx context.Context, tx *gorm.DB,accountID string) error{
	if err := tx.Table(models.Account{}.TableName()).
		Where("id = ?", accountID).
		Update("account_status",false).Error; err != nil {
		return err
	}
	return nil
}