package auth

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	minPasswordLength = 8
)

func (a *User) Validate() error {
	if a.Username == "" {
		return errors.New("username is empty")
	}
	if a.Email == "" {
		return errors.New("email is empty")
	}
	if a.Password == "" {
		return errors.New("password is empty")
	}
	if len(a.Password) < minPasswordLength {
		return errors.New("password is too short")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
