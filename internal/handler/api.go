package handler

import (
	"balance/internal/apperror"
	"balance/internal/model"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var a *model.Account

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	if err := c.ShouldBindJSON(&a); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	a.UserID = userID

	// проверка счета пользователя на дубликат
	existsAccount, err := h.Service.ExistsAccount(a.Number)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}
	if existsAccount {
		c.JSON(http.StatusBadRequest, apperror.ErrExistsAccount)
		return
	}

	// регистрация нового счета пользователя
	err = h.Service.CreateAccount(a)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusCreated, "Adding new account was successful")
}

func (h *Handler) GetAccounts(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	accounts, err := h.Service.GetAccounts(userID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) GetAccountById(c *gin.Context) {
	id := c.Query("id")

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	account, err := h.Service.GetAccountById(userID, id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	var a *model.Account

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	a.UserID = userID

	err := c.ShouldBindJSON(&a)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	err = h.Service.UpdateAccount(a)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, "Updating account was successful")
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var tr *model.Transaction

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	err := c.ShouldBindJSON(&tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}
	ok = tr.Type == "expense" || tr.Type == "income"
	if !ok {
		h.Logger.Error(fmt.Errorf("incorrect transation type"))
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	account, err := h.Service.GetAccountById(userID, tr.AccountID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}
	if account.ID == "" {
		h.Logger.Errorf("this user doesn't have an account with id %s", tr.AccountID)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	err = h.Service.CreateTransaction(tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	if tr.Type == "expense" {
		account.Balance -= tr.Amount
	} else if tr.Type == "income" {
		account.Balance += tr.Amount
	}

	err = h.Service.UpdateAccount(&account)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusCreated, "the transaction is saved")
}

func (h *Handler) GetTransactions(c *gin.Context) {
	var tr []model.Transaction

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	accounts, err := h.Service.GetAccounts(userID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	for _, account := range accounts {
		transactions, err := h.Service.GetTransactions(account.ID)
		if err != nil {
			h.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
			return
		}

		tr = append(tr, transactions...)
	}

	c.JSON(http.StatusOK, tr)
}

func (h *Handler) GetTransactionById(c *gin.Context) {
	id := c.Query("id")

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	transaction, err := h.Service.GetTransactionById(id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	account, err := h.Service.GetAccountInfoById(transaction.AccountID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}
	if account.UserID != userId {
		h.Logger.Error("this transaction doesn't belong to user")
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) GetReports(c *gin.Context) {
	var rep *model.Report

	userId, ok := c.Get("user_id")
	if !ok {
		h.Logger.Error("cannot to get user ID from token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userID := userId.(string)

	err := c.ShouldBindJSON(&rep)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}
	if (rep.Type != "") && (rep.Type != "expense") && (rep.Type != "income") {
		h.Logger.Error(fmt.Errorf("incorrect transation type"))
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	from := time.Time{}
	to := time.Time{}
	if rep.DateFrom != "" {
		from, err = time.Parse("02-01-2006", rep.DateFrom)
		if err != nil {
			h.Logger.Error(err)
			c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
			return
		}
	}
	if rep.DateFrom != "" {
		to, err = time.Parse("02-01-2006", rep.DateTo)
		if err != nil {
			h.Logger.Error(err)
			c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
			return
		}
	}

	rep.From = from
	rep.To = to

	reports, err := h.Service.GetReports(userID, rep)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	// Сохраняем файл в буфер
	buffer := new(bytes.Buffer)
	err = reports.Write(buffer)
	if err != nil {
		h.Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Отправляем файл клиенту
	c.Header("Content-Disposition", "attachment; filename=example.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
