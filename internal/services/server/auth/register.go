package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// RegisterUser - регистрация пользователя.
func (s *ServiceAuth) RegisterUser(ctx context.Context, login, password string) (int, error) {

	err := s.storage.CheckUser(ctx, login)
	if err != nil {
		return -1, err
	}

	hashPassword := s.hashPassword(password)

	err = s.storage.SaveUser(ctx, login, hashPassword)
	if err != nil {
		return -1, err
	}

	uid, err := s.storage.GetUserIDByLogin(ctx, login)
	if err != nil {
		return -1, err
	}

	return uid, nil
}

// hashPassword - хеширование пароля.
func (s *ServiceAuth) hashPassword(password string) string {
	var passwordBytes = []byte(password)
	var sha512Hashes = sha256.New()

	passwordBytes = append(passwordBytes, s.passwordSalt...)

	sha512Hashes.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hashes.Sum(nil)
	var hashedPasswordHax = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHax
}
