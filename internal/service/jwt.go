package service

import (
	"balance/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"os"
	"time"
)

func (s *Service) CreateToken(userID string) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	// Creating Access Token
	err := os.Setenv("ACCESS_SECRET", "secret")
	if err != nil {
		return nil, err
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating Refresh Token
	err = os.Setenv("REFRESH_SECRET", "secret")
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *Service) CreateAuth(userID string, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := s.Redis.Set(td.AccessUuid, userID, at.Sub(now)).Err()
	if errAccess != nil {
		s.Logger.Error(errAccess)
		return errAccess
	}
	errRefresh := s.Redis.Set(td.RefreshUuid, userID, rt.Sub(now)).Err()
	if errRefresh != nil {
		s.Logger.Error(errRefresh)
		return errRefresh
	}

	return nil
}

func (s *Service) DeleteAuth(givenUuid string) (int64, error) {

	deleted, err := s.Redis.Del(givenUuid).Result()
	if err != nil {
		s.Logger.Error(err)
		return 0, err
	}

	return deleted, nil
}
