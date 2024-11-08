package auth

import (
	"encoding/json"
	"github.com/Fyefhqdishka/web-project/pkg/jwt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type RepositoryAuth interface {
	RegistrationUser(user User) error
	LoginUser(Login, Username string) (string, error)
}

type ControllerAuth struct {
	repo   RepositoryAuth
	logger *slog.Logger
}

func NewControllerAuth(repo RepositoryAuth, logger *slog.Logger) *ControllerAuth {
	return &ControllerAuth{
		repo:   repo,
		logger: logger,
	}
}
func (c *ControllerAuth) Register(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug(
		"Register",
		"начала обработки регистрации пользователя",
	)

	//files := []string{
	//	"./internal/ui/auth.html",
	//}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		c.logger.Error(
			"Register",
			"ошибка при декодировании JSON",
			"err", err,
		)
		return
	}

	if err = user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.logger.Error(
			"Register",
			"ошибка валидации данных пользователя",
			"err", err,
		)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.logger.Error(
			"Register",
			"ошибка при хэшировании пароля",
			"err:", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	if err = c.repo.RegistrationUser(user); err != nil {
		c.logger.Error(
			"Register",
			"ошибка при сохранении пользователя в репозитории",
			"err", err,
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	c.logger.Debug(
		"Register",
		"User successfully registered",
		"username", user.Username,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (c *ControllerAuth) Login(w http.ResponseWriter, r *http.Request) {
	c.logger.Info(
		"login",
		"начало обработки Авторизации",
	)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Error(
			"Login",
			"ошибка чтения тела запроса",
			"error", err,
		)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var user *User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		c.logger.Error(
			"Login",
			"ошибка декодирования JSON",
			"error", err,
		)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	c.logger.Debug(
		"Login",
		"Перед вызовом LoginUser",
		"username", user.Username,
	)

	userID, err := c.repo.LoginUser(user.Username, user.Password)
	if err != nil {
		c.logger.Error(
			"Login",
			"ошибка при аутентификации пользователя",
			"username", user.Username,
			"error", err,
		)

		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	tokenStr, err := jwt.GenerateToken(userID)
	if err != nil {
		c.logger.Error(
			"Login",
			"Ошибка генерации токена",
			"error", err,
		)

		http.Error(w, "Не удалось сгенерировать токен", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  time.Now().Add(36 * time.Hour),
		HttpOnly: true,
	})

	c.logger.Info(
		"Login",
		"Пользователь успешно аутентифицирован",
		"username", user.Username,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User authenticate successfully",
		"token":   tokenStr,
	})
}
