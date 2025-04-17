package utils

import (
	"time"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Error возвращает карту с ошибкой для JSON-ответа
func Error(msg string) fiber.Map {
	return fiber.Map{"error": msg}
}

// CheckPassword сравнивает обычный пароль с хешем
// Возвращает true, если пароль верный
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT создаёт и подписывает JWT-токен для пользователя
func GenerateJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET") // Чтение секрета из .env

	// Создаём токен с user_id и сроком действия
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Токен живёт 3 дня
	})

	return token.SignedString([]byte(secret))
}