package auth

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		db,
		logger,
	}
}

func (m *Repository) RegistrationUser(user User) error {
	stmt, err := m.DB.Prepare("INSERT INTO users (name, username, email, password) VALUES ($1,$2,$3,$4)")
	if err != nil {
		m.Logger.Error(
			"RegistrationUser",
			"error preparing statement",
			"err:", err,
		)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		m.Logger.Error(
			"RegistrationUser",
			"error executing statement",
			"err:", err,
		)

		return err
	}

	m.Logger.Debug(
		"RegistrationUser",
		"registration user finished",
		"username", user.Username,
	)

	return nil
}

func (m *Repository) LoginUser(username, password string) (string, error) {
	m.Logger.Debug(
		"LoginUser",
		"starting login user",
	)

	var UserID string
	var passwordHash string

	stmt := `SELECT id, password FROM users WHERE username = $1`
	err := m.DB.QueryRow(stmt, username).Scan(&UserID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			m.Logger.Warn(
				"LoginUser",
				"пользователь не найден",
				"username", username,
			)

			return "", fmt.Errorf("пользователь не найден")
		}
		m.Logger.Error(
			"LoginUser",
			"ошибка выполнения SQL-запроса при логине",
			"err", err,
		)

		return "", err
	}

	if !CheckPasswordHash(password, passwordHash) {
		m.Logger.Warn(
			"LoginUser",
			"неверный пароль для пользователя",
			"username", username,
		)

		return "", fmt.Errorf("неверный пароль")
	}

	m.Logger.Info(
		"пользователь успешно аутентифицирован",
		"username", username,
	)

	return UserID, nil
}
