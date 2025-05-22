# 🎥 🍿  
# Movie Recommendation API

Сервис для управления фильмов с аутентификации пользователей, оценки кино, лайков и выдачи персональных рекомендаций. Разработан на Go с использованием REST API, PostgreSQL и интеграцией с TMDB. 

---

## Описание

Сервис позволяет:

1. **Регистрироваться и логиниться** — безопасная регистрация с хешированием паролей и JWT.
2. **Оценивать фильмы** — ставить оценку от 1 до 10.
3. **Лайкать фильмы** — сохранять любимые фильмы в избранное.
4. **Получать рекомендации** — (в будущем) — на основе лайков и оценок.
5. **Искать фильмы** — поиск по базе TMDB.

---

## Основные функции

### Поддерживаемые методы:

- **POST** `/register`
- **POST** `/login`
- **GET** `/me`
- **POST** `/movies/like`
- **GET** `/movies/liked`
- **POST** `/movies/rate`
- **GET** `/movies/rated`
- **DELETE** `/movies/unlike`, `/movies/unrate`
- **GET** `/search?query=...`

---

## Стек технологий

- **Go + Fiber**
- **PostgreSQL + GORM**
- **JWT** — аутентификация
- **Docker** и **Docker Compose**
- **TMDB API** — внешний источник фильмов

---

## Установка и запуск

```bash
git clone https://github.com/ot4u/movie_service.git
cd movie_service
make run
```

### Переменная окружения .env

```.env
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/moviedb?sslmode=disable
JWT_SECRET=supersecret123
TMDB_API_KEY=your_tmdb_key_here
```

## Тестирование

```bash
make test
```

## Архитектура

Проект разделен на слои:
- **handlers** — REST-эндпоинты
- **models** — структуры БД 
- **middleware** — JWT защита
- **services** — внешние API (TMDB)
- **database** — подключение и миграции

Структура проекта:
```
├── cmd
│   └── main.go
├── internal
│   ├── handlers
│   ├── middleware
│   ├── models
│   ├── services
│   └── database
├── pkg
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```
## Возможные улучшения

- Redis и Kafka
- Роли и права
- Веб-интерфейс

---
> **“Хороший код — это не когда нечего добавить, а когда нечего убрать.”**