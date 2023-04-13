package repository

import (
	"balance/internal/apperror"
	"balance/internal/model"
	"balance/logging"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	Connection *gorm.DB
	Logger     *logging.Logger
}

func NewRepository(connection *gorm.DB, logger *logging.Logger) *Repository {
	return &Repository{
		Connection: connection,
		Logger:     logger,
	}
}

func (r *Repository) ExistsUser(username string) (bool, error) {
	var u model.User
	err := r.Connection.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if u == (model.User{}) {
		return false, nil
	}
	return true, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) (userID string, err error) {
	err = r.Connection.WithContext(ctx).Create(&user).Error
	if err != nil {
		r.Logger.Errorf("failed to create user due error: %v", err)
		return "", err
	}

	return user.ID, nil
}

func (r *Repository) CheckUser(user *model.User) (u *model.User, err error) {
	if tx := r.Connection.Where("username = ?", user.Username).Find(&u); tx.Error != nil {
		r.Logger.Errorf("failed to user due error: %v", tx.Error)
		return u, tx.Error
	}
	if u == (&model.User{}) {
		r.Logger.Error(apperror.ErrNotFound)
		return u, apperror.ErrNotFound
	}

	return u, nil
}

func (r *Repository) ExistsAccount(number string) (bool, error) {
	var a *model.Account

	err := r.Connection.Where("number = ?", number).Find(&a).Error
	if err != nil {
		r.Logger.Error(err)
		return false, err
	}
	if a.Number != "" {
		return true, nil
	}

	return false, nil
}

func (r *Repository) CreateAccount(account *model.Account) error {
	err := r.Connection.Omit("created", "updated", "deleted").Create(&account).Error
	if err != nil {
		r.Logger.Error(err)
		return err
	}

	return nil
}

func (r *Repository) GetAccounts(userID string) (accounts []model.Account, err error) {
	err = r.Connection.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return accounts, nil
}

func (r *Repository) GetAccountById(userId, id string) (account model.Account, err error) {
	err = r.Connection.Where("user_id = ? and id = ?", userId, id).Find(&account).Error
	if err != nil {
		r.Logger.Error(err)
		return model.Account{}, err
	}

	return account, nil
}

func (r *Repository) UpdateAccount(account *model.Account) error {
	err := r.Connection.Save(account).Error
	if err != nil {
		r.Logger.Error(err)
		return err
	}

	return nil
}

func (r *Repository) CreateTransaction(tr *model.Transaction) error {
	err := r.Connection.Omit("created", "updated", "deleted").Create(&tr).Error
	if err != nil {
		r.Logger.Error(err)
		return err
	}

	return nil
}

func (r *Repository) GetUserIdByAccountID(accountID string) (string, error) {
	var a *model.Account

	err := r.Connection.Where("id = ?", accountID).Find(&a).Error
	if err != nil {
		r.Logger.Error(err)
		return "", err
	}
	if a.UserID == "" {
		return "", fmt.Errorf("does not exists account with id: %s", accountID)
	}

	return a.UserID, nil
}

func (r *Repository) GetTransactions(userID string) (tr []model.Transaction, err error) {
	err = r.Connection.Where("account_id = ?", userID).Find(&tr).Error
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return tr, nil
}

func (r *Repository) GetTransactionById(id string) (tr model.Transaction, err error) {
	err = r.Connection.Where("id = ?", id).Find(&tr).Error
	if err != nil {
		r.Logger.Error(err)
		return model.Transaction{}, err
	}

	return tr, nil
}

func (r *Repository) GetReports(rep *model.Report) (tr []model.Transaction, err error) {
	query := r.Connection
	if rep.Type != "" {
		query = query.Where("type = ?", rep.Type)
	}
	if rep.From != (time.Time{}) {
		query = query.Where("created >= ?", rep.From)
	}
	if rep.To != (time.Time{}) {
		query = query.Where("created <= ?", rep.To)
	}

	page := 1
	limit := 0

	if rep.Page > 0 {
		page = rep.Page
	}
	if rep.Limit > 0 {
		limit = rep.Limit
	}

	if rep.Page > 0 {
		query = query.Limit(limit).Offset((page - 1) * limit)
	}

	err = query.Find(&tr).Error
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return tr, nil
}

func (r *Repository) GetInfoByUserId(userID string) (u *model.User, err error) {
	err = r.Connection.Where("id = ?", userID).Find(&u).Error
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return u, nil
}

func (r *Repository) GetAccountInfoById(accountID string) (a *model.Account, err error) {
	err = r.Connection.Where("id = ?", accountID).Find(&a).Error
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return a, nil
}
