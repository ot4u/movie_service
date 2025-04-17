package handlers

import (
	"movie_service/internal/database"
	"movie_service/internal/models"
	"movie_service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var input LoginInput

	// 1. Парсим тело запроса в структуру LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Неверный формат запроса"))
	}

	// 2. Ищем пользователя по email в базе данных
	user, err := findUserByEmail(input.Email)
	if err != nil {
		// Не говорим "email не найден", чтобы не палить существующих пользователей
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Неверный email или пароль"))
	}

	// 3. Проверяем пароль: сравниваем с хешем из базы
	if !utils.CheckPassword(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Неверный email или пароль"))
	}

	// 4. Генерируем JWT токен с user_id
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Ошибка при создании токена"))
	}

	// 5. Возвращаем токен клиенту
	return c.JSON(fiber.Map{"token": token})
}

func findUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}