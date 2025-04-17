package handlers

import (
	"movie_service/internal/database"
	"movie_service/internal/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	var status = fiber.StatusOK
	var response fiber.Map

	// 1. Парсим JSON
	err := c.BodyParser(&input)

	if err != nil {
		status = fiber.StatusBadRequest
		response = fiber.Map{"error": "Неверный формат запроса"}
	} else {
		// 2. Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			status = fiber.StatusInternalServerError
			response = fiber.Map{"error": "Ошибка при хешировании пароля"}
		} else {

			// 3. Создаём пользователя
				user := models.User{
				Username: input.Username,
				Email:    input.Email,
				Password: string(hashedPassword),
			}

			// 4. Сохраняем в БД
			result := database.DB.Create(&user)
			if result.Error != nil {
				status = fiber.StatusInternalServerError
				response = fiber.Map{"error": "Не удалось создать пользователя"}
			} else {

				// 5. Всё ок
				response = fiber.Map{"message": "Регистрация прошла успешно"}
			}
		}
	}

	return c.Status(status).JSON(response)
}