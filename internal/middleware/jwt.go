package middleware

import (
	"movie_service/internal/database"
	"movie_service/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Читаем заголовок Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Нет токена в заголовке",
			})
		}

		// 2. Вырезаем токен из строки "Bearer <токен>"
		tokenStr := authHeader[len("Bearer "):]

		// 3. Получаем секрет из .env
		secret := []byte(os.Getenv("JWT_SECRET"))

		// 4. Парсим и проверяем токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Неверный токен",
			})
		}

		// 5. Извлекаем user_id из claims
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		// 6. Загружаем пользователя из БД
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Пользователь не найден",
			})
		}

		// 7. Кладём пользователя в контекст запроса
		c.Locals("user", user)

		// 8. Переходим к следующему хендлеру
		return c.Next()
	}
}