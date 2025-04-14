package main

import (
	"log"
	"movie_service/internal/database"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"movie_service/internal/handlers"
)

func main() {
	app := fiber.New() // Создаём сервер

	// Роут для проверки работы
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Movie API is running 🚀")
	})

	//Подключение файла окружения
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не обнаружен")
	}

	//Чтение данных из файла окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // значение по умолчанию, если PORT не задан
	}

	//Подключение к бд
	database.Connect()
	
	//Регистрация пользователя
	app.Post("/register", handlers.Register)

	//Запуск сервера
	log.Fatal(app.Listen(":" + port))

}