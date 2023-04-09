package handler

import (
	"balance/internal/service"
	"balance/logging"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Engine  *gin.Engine
	Logger  *logging.Logger
	Service *service.Service
}

func NewHandler(engine *gin.Engine, logger *logging.Logger, service *service.Service) *Handler {
	return &Handler{
		Engine:  engine,
		Logger:  logger,
		Service: service,
	}
}

func (h *Handler) Init() {
	generalRout := h.Engine.Group("v1")

	auth := generalRout.Group("/auth")
	{
		auth.POST("/create", h.CreateUser)
		auth.POST("/login", h.Login)
		auth.POST("/token/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	api := generalRout.Group("/api")
	api.Use(h.TokenAuthMiddleware())
	{
		api.POST("/account", h.CreateAccount)
		api.GET("/account", h.GetAccounts)
		api.GET("/account/:id", h.GetAccountById)
		api.PUT("/account", h.UpdateAccount)
		api.POST("/transaction", h.CreateTransaction)
		api.GET("/transaction", h.GetTransactions)
		api.GET("/transaction/:id", h.GetTransactionById)
		api.GET("/excel-reports", h.GetReports)
	}
}
