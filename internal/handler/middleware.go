package handler

import (
	"balance/internal/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func (h *Handler) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := h.TokenValid(c.Request)
		if err != nil {
			h.Logger.Error(err)
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", userID)

		c.Next()
	}
}

func (h *Handler) TokenValid(r *http.Request) (string, error) {
	token, err := h.VerifyToken(r)
	if err != nil {
		h.Logger.Error(err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		h.Logger.Error(err)
		return "", err
	}

	userID := claims["user_id"].(string)

	return userID, nil
}

func (h *Handler) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := h.ExtractToken(r)

	err := os.Setenv("ACCESS_SECRET", "secret")
	if err != nil {
		h.Logger.Error(err)
		return nil, fmt.Errorf("cannot set data")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method confirm to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.Logger.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
		// return []byte("secret"), nil
	})
	if err != nil {
		h.Logger.Error(err)
		return nil, err
	}

	return token, nil
}

func (h *Handler) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("token")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return strArr[0]
}

func (h *Handler) ExtractTokenMetaData(r *http.Request) (*model.AccessDetails, error) {
	token, err := h.VerifyToken(r)
	if err != nil {
		h.Logger.Error(err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			h.Logger.Error(err)
			return nil, err
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			// TODO
		}
		return &model.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}
