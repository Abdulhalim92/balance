package handler

import (
	"balance/internal/apperror"
	"balance/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var a *model.Account

	if err := c.ShouldBindJSON(&a); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

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
	var a *model.Account

	err := c.ShouldBindJSON(&a)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	accounts, err := h.Service.GetAccounts(a.UserID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) GetAccountById(c *gin.Context) {
	id := c.Query("id")

	var a *model.Account

	err := c.ShouldBindJSON(&a)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	account, err := h.Service.GetAccountById(a.UserID, id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	var a *model.Account

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

	err := c.ShouldBindJSON(&tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}
	if (tr.Type != "expense") || (tr.Type != "income") {
		h.Logger.Error(fmt.Errorf("incorrect transation type"))
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	err = h.Service.CreateTransaction(tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	userID, err := h.Service.Repository.GetUserIdByAccountID(tr.AccountID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	account, err := h.Service.GetAccountById(userID, tr.AccountID)
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
	var tr *model.Transaction

	err := c.ShouldBindJSON(&tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	//userID, err := h.Service.Repository.GetUserIdByAccountID(tr.AccountID)
	//if err != nil {
	//	h.Logger.Error(err)
	//	c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
	//	return
	//}

	transactions, err := h.Service.GetTransactions(tr.AccountID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) GetTransactionById(c *gin.Context) {
	id := c.Query("id")

	var tr *model.Transaction

	err := c.ShouldBindJSON(&tr)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	userID, err := h.Service.Repository.GetUserIdByAccountID(tr.AccountID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	transaction, err := h.Service.GetTransactionById(id)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) GetReports(c *gin.Context) {
	var rep *model.Report

	err := c.ShouldBindJSON(&rep)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}
	if (rep.Type != "expense") || (rep.Type != "income") {
		h.Logger.Error(fmt.Errorf("incorrect transation type"))
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	from, err := time.Parse("02-01-2006", rep.DateFrom)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	to, err := time.Parse("02-01-2006", rep.DateTo)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	rep.From = from
	rep.To = to

	h.Service.GetReports(rep)
}
