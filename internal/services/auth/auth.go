package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Tokens struct {
	AccessToken string
}

type AccessTokenClaims struct {
	Login string
	jwt.RegisteredClaims
}

var (
	errNotFoundPassword = errors.New("not found password")
	errWrongPassword    = errors.New("wrong password")
)

// Auth - авторизация пользователя.
func (s *ServiceAuth) Auth(login, password string) (Tokens, error) {
	passwordHash, ok := s.storage.CheckPassword(login)
	if !ok {
		return Tokens{}, errNotFoundPassword
	}

	if !s.doPasswordMatch(passwordHash, password) {
		return Tokens{}, errWrongPassword
	}

	token, err := s.generateToken(login)
	if err != nil {
		return Tokens{}, err
	}

	return token, nil
}

// doPasswordMatch - сравнение хешированных паролей.
func (s *ServiceAuth) doPasswordMatch(hashedPassword, password string) bool {
	var currPasswordHash = s.hashPassword(password)

	return hashedPassword == currPasswordHash
}

// generateToken - генерация токена.
func (s *ServiceAuth) generateToken(login string) (Tokens, error) {
	accessToken, err := s.generateAccessToken(login)
	if err != nil {
		return Tokens{}, err
	}

	// сохраняем accessToken в базу users
	err = s.storage.SaveTableUserAndUpdateToken(login, accessToken)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken: accessToken,
	}, nil
}

// generateAccessToken - генерация accessToken.
func (s *ServiceAuth) generateAccessToken(login string) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.tokenSalt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
