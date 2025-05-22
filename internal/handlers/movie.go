package handlers

import (
	"movie_service/internal/database"
	"movie_service/internal/models"
	"movie_service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LikeMovieInput struct {
	TMDB_ID     int    `json:"tmdb_id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

type RateMovieInput struct {
	TMDB_ID     int    `json:"tmdb_id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
	Score       int    `json:"score"`
}

type TMDB_ID_Input struct {
	TMDB_ID int `json:"tmdb_id"`
}

func LikeMovie(c *fiber.Ctx) error {
	var input LikeMovieInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	// Получаем текущего пользователя
	user := c.Locals("user").(models.User)

	// Ищем или создаём фильм
	var movie models.Movie
	db := database.DB
	err := db.Where("tmdb_id = ?", input.TMDB_ID).First(&movie).Error
	if err != nil {
		movie = models.Movie{
			TMDB_ID:     input.TMDB_ID,
			Title:       input.Title,
			PosterPath:  input.PosterPath,
			ReleaseDate: input.ReleaseDate,
		}
		db.Create(&movie)
	}

	// Проверка: лайкал ли уже
	var existingLike models.UserLike
	err = db.Where("user_id = ? AND movie_id = ?", user.ID, movie.ID).First(&existingLike).Error
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Фильм уже добавлен в избранное"})
	}

	// Сохраняем лайк
	like := models.UserLike{
		UserID:  user.ID,
		MovieID: movie.ID,
	}
	db.Create(&like)

	return c.JSON(fiber.Map{"message": "Фильм добавлен в избранное"})
}

func GetLikedMovies(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var likedMovies []models.Movie

	err := database.DB.
		Joins("JOIN user_likes ON user_likes.movie_id = movies.id").
		Where("user_likes.user_id = ?", user.ID).
		Find(&likedMovies).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить избранные фильмы",
		})
	}

	// Преобразуем в DTO
	var response []models.MovieResponse
	for _, m := range likedMovies {
		response = append(response, models.MovieResponse{
			TMDB_ID:     m.TMDB_ID,
			Title:       m.Title,
			PosterPath:  m.PosterPath,
			ReleaseDate: m.ReleaseDate,
		})
	}

	return c.JSON(response)
}

func RateMovie(c *fiber.Ctx) error {
	var input RateMovieInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	if input.Score < 1 || input.Score > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Оценка должна быть от 1 до 10"})
	}

	user := c.Locals("user").(models.User)

	// Найти или создать фильм
	var movie models.Movie
	err := database.DB.Where("tmdb_id = ?", input.TMDB_ID).First(&movie).Error
	if err != nil {
		movie = models.Movie{
			TMDB_ID:     input.TMDB_ID,
			Title:       input.Title,
			PosterPath:  input.PosterPath,
			ReleaseDate: input.ReleaseDate,
		}
		database.DB.Create(&movie)
	}

	// Проверяем: уже есть рейтинг?
	var rating models.Rating
	err = database.DB.
		Where("user_id = ? AND movie_id = ?", user.ID, movie.ID).
		First(&rating).Error

	if err == nil {
		// Обновляем рейтинг
		rating.Score = input.Score
		database.DB.Save(&rating)
	} else {
		// Создаём новый рейтинг
		rating = models.Rating{
			UserID:  user.ID,
			MovieID: movie.ID,
			Score:   input.Score,
		}
		database.DB.Create(&rating)
	}

	return c.JSON(fiber.Map{"message": "Оценка сохранена"})
}

func GetRatedMovies(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var ratings []models.Rating

	// Загружаем все оценки пользователя с данными фильма
	err := database.DB.
		Preload("Movie").
		Where("user_id = ?", user.ID).
		Find(&ratings).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось загрузить оценки",
		})
	}

	// Преобразуем в DTO
	var response []models.RatedMovieResponse
	for _, r := range ratings {
		response = append(response, models.RatedMovieResponse{
			TMDB_ID:     r.Movie.TMDB_ID,
			Title:       r.Movie.Title,
			PosterPath:  r.Movie.PosterPath,
			ReleaseDate: r.Movie.ReleaseDate,
			Score:       r.Score,
		})
	}

	return c.JSON(response)
}

func UnlikeMovie(c *fiber.Ctx) error {
	var input TMDB_ID_Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	user := c.Locals("user").(models.User)

	var movie models.Movie
	err := database.DB.Where("tmdb_id = ?", input.TMDB_ID).First(&movie).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Фильм не найден"})
	}

	err = database.DB.
		Where("user_id = ? AND movie_id = ?", user.ID, movie.ID).
		Delete(&models.UserLike{}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления лайка"})
	}

	return c.JSON(fiber.Map{"message": "Лайк удалён"})
}

func UnrateMovie(c *fiber.Ctx) error {
	var input TMDB_ID_Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	user := c.Locals("user").(models.User)

	var movie models.Movie
	err := database.DB.Where("tmdb_id = ?", input.TMDB_ID).First(&movie).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Фильм не найден"})
	}

	err = database.DB.
		Where("user_id = ? AND movie_id = ?", user.ID, movie.ID).
		Delete(&models.Rating{}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления оценки"})
	}

	return c.JSON(fiber.Map{"message": "Оценка удалена"})
}

func SearchMovies(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Параметр 'query' обязателен"})
	}

	results, err := services.SearchMovies(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при запросе к TMDB"})
	}

	return c.JSON(results)
}