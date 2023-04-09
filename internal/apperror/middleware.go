package apperror

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type appHandler func(c *gin.Context) error

func Middleware(handler appHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var appError *AppError
		err := handler(c)
		if err != nil {
			if errors.As(err, &appError) {
				if errors.Is(err, ErrNotFound) {
					c.JSON(http.StatusNotFound, ErrNotFound.Marshal())
					return
				}
				err = err.(*AppError)
				c.JSON(http.StatusBadRequest, appError.Marshal())
				return
			}

			c.JSON(http.StatusTeapot, systemError(err).Marshal())
		}
	}
}
