package service

import (
	"balance/internal/model"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func (s *Service) ExistsAccount(number string) (bool, error) {
	existsAccount, err := s.Repository.ExistsAccount(number)
	if err != nil {
		s.Logger.Error(err)
		return false, err
	}

	return existsAccount, nil
}

func (s *Service) CreateAccount(account *model.Account) error {
	err := s.Repository.CreateAccount(account)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	return nil
}

func (s *Service) GetAccounts(userID string) ([]model.Account, error) {
	accounts, err := s.Repository.GetAccounts(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return accounts, nil
}

func (s *Service) GetAccountById(userID, id string) (model.Account, error) {
	account, err := s.Repository.GetAccountById(userID, id)
	if err != nil {
		s.Logger.Error(err)
		return model.Account{}, err
	}

	return account, nil
}

func (s *Service) UpdateAccount(account *model.Account) error {
	err := s.Repository.UpdateAccount(account)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	return nil
}

func (s *Service) CreateTransaction(tr *model.Transaction) error {
	err := s.Repository.CreateTransaction(tr)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	return nil
}

func (s *Service) GetTransactions(accountID string) ([]model.Transaction, error) {
	tr, err := s.Repository.GetTransactions(accountID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return tr, nil
}

func (s *Service) GetTransactionById(id string) (model.Transaction, error) {
	tr, err := s.Repository.GetTransactionById(id)
	if err != nil {
		s.Logger.Error(err)
		return model.Transaction{}, err
	}

	return tr, nil
}

func (s *Service) GetReports(userID string, rep *model.Report) (*excelize.File, error) {
	tr, err := s.Repository.GetReports(rep)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	excelReports, err := s.GetExcelReports(userID, tr)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return excelReports, nil
}

func (s *Service) GetExcelReports(userID string, tr []model.Transaction) (*excelize.File, error) {
	excelFile := excelize.NewFile()

	sheet, err := excelFile.NewSheet("Отчёт")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "A1", "Имя пользователя")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	err = excelFile.SetCellValue("Отчёт", "B1", "Название счёта")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	err = excelFile.SetCellValue("Отчёт", "C1", "Тип операции")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	err = excelFile.SetCellValue("Отчёт", "D1", "Сумма")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	err = excelFile.SetCellValue("Отчёт", "E1", "Дата совершения операции")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	u, err := s.GetUserInfoById(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	for i, transaction := range tr {
		i += 2

		account, err := s.GetAccountInfoById(transaction.AccountID)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "A"+strconv.Itoa(i), u.Username)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		err = excelFile.SetCellValue("Отчёт", "B"+strconv.Itoa(i), account.Number)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		err = excelFile.SetCellValue("Отчёт", "C"+strconv.Itoa(i), transaction.Type)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		err = excelFile.SetCellValue("Отчёт", "D"+strconv.Itoa(i), transaction.Amount)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		err = excelFile.SetCellValue("Отчёт", "E"+strconv.Itoa(i), transaction.Created)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

	}

	excelFile.SetActiveSheet(sheet)

	err = excelFile.SaveAs("report.xlsx")
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return excelFile, nil
}

func (s *Service) GetUserInfoById(userID string) (*model.User, error) {
	u, err := s.Repository.GetInfoByUserId(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return u, nil
}

func (s *Service) GetAccountInfoById(accountID string) (*model.Account, error) {
	account, err := s.Repository.GetAccountInfoById(accountID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return account, nil
}
