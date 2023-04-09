package service

import (
	"balance/internal/apperror"
	"balance/internal/model"
	"context"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func (s *Service) ValidateUser(user *model.User) error {
	if len(user.Username) > 20 || len(user.Username) < 3 {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if len(user.Password) > 20 || len(user.Password) < 6 {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "_") || strings.Contains(user.Password, "-") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "@") || strings.Contains(user.Password, "#") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "$") || strings.Contains(user.Password, "%") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "&") || strings.Contains(user.Password, "*") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "(") || strings.Contains(user.Password, ")") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, ":") || strings.Contains(user.Password, ".") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "/") || strings.Contains(user.Password, `\`) {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, ",") || strings.Contains(user.Password, ";") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "?") || strings.Contains(user.Password, `"`) {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "!") || strings.Contains(user.Password, "~") {
		s.Logger.Error(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	return nil
}

func (s *Service) ExistsUser(username string) (bool, error) {
	existsUser, err := s.Repository.ExistsUser(username)
	if err != nil {
		s.Logger.Error(err)
		return false, err
	}
	if existsUser {
		s.Logger.Infoln("this username is registered")
		return true, nil
	}

	return false, nil
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Errorf("falied to generate hash from password due error: %v", err)
		return "", err
	}

	user.Password = string(hashPassword)

	userID, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		s.Logger.Errorf("failed to create user due error: %v", err)
		return "", nil
	}

	return userID, nil
}

func (s *Service) CheckUser(user *model.User) (string, error) {
	u, err := s.Repository.CheckUser(user)
	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	return u.ID, nil
}
