package handler

import (
	"balance/internal/apperror"
	"balance/internal/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var u *model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()

	err = h.Service.ValidateUser(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrInvalid)
		return
	}

	existsUser, err := h.Service.ExistsUser(u.Username)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}
	if existsUser {
		c.JSON(http.StatusBadRequest, apperror.ErrRegistered)
		return
	}

	userID, err := h.Service.CreateUser(ctx, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"user_id": userID,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var u *model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	userID, err := h.Service.CheckUser(u)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusUnauthorized, apperror.ErrUnauthorized)
		return
	}

	ts, err := h.Service.CreateToken(userID)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	saveErr := h.Service.CreateAuth(userID, ts)
	if saveErr != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {
	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}

	refreshToken := mapToken["refresh_token"]

	// verify the token
	err := os.Setenv("REFRESH_SECRET", "secret")
	if err != nil {
		h.Logger.Error(err)
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		h.Logger.Error(apperror.ErrExpiredRefresh)
		c.JSON(http.StatusUnauthorized, apperror.ErrExpiredRefresh)
		return
	}

	// is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		h.Logger.Error(apperror.ErrInvalidToken)
		c.JSON(http.StatusUnauthorized, apperror.ErrInvalidToken)
		return
	}

	// since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			h.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
			return
		}
		//userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		//if err != nil {
		//	c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		//	return
		//}

		userID := claims["user_id"].(string)

		// delete the previous Refresh Token
		deleted, delErr := h.Service.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			h.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		}

		// create new pairs of refresh and access tokens
		ts, createErr := h.Service.CreateToken(userID)
		if createErr != nil {
			h.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, apperror.ErrExpiredToken)
	}
}

func (h *Handler) Logout(c *gin.Context) {
	au, err := h.ExtractTokenMetaData(c.Request)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, apperror.ErrBadRequest)
		return
	}
	deleted, delErr := h.Service.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		h.Logger.Error(delErr)
		c.JSON(http.StatusInternalServerError, apperror.ErrInternalServer)
		return
	}
	c.JSON(http.StatusOK, "Successfully loddeg out")
}
