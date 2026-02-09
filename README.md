<div align="center">

# My Blog Backend

[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-05122A?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-blue?style=for-the-badge)](https://gorm.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io/)
[![Swagger](https://img.shields.io/badge/-Swagger-%23C0E800?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)

**Backend-сервис для платформы блога, предоставляющий API для управления пользователями, статьями и категориями.**

</div>

## 🚀 Быстрый старт

1.  **Клонируйте репозиторий:**
    ```bash
    git clone git@github.com:SkyGreenxd/my_blog_backend.git
    cd my_blog_backend
    ```

2.  **Настройте переменные окружения:**
    Скопируйте пример конфига и отредактируйте его при необходимости:
    ```bash
    cp .env.example .env
    ```

3.  **Запустите проект:**
    ```bash
    docker-compose up --build
    ```
    После этого API будет доступно по адресу `http://localhost:8080`.
---

## 📖 API Documentation
Swagger: http://localhost:8080/api/v1/swagger/index.html

## 🏛 Структура проекта

- `cmd/app`: Точка входа в приложение.
- `db/migrations`: SQL файлы для миграции базы данных.
- `internal/config`: Логика загрузки конфигурации.
- `internal/delivery/v1`: Обработка HTTP запросов (Handlers & Middleware).
- `internal/repository`: Работа с базой данных (GORM).
- `internal/server`: Кастомный HTTP сервер. 
- `internal/domain`: Доменные сущности.
- `internal/usecase`: Бизнес-логика приложения.
- `pkg/`: Вспомогательные утилиты и пакеты общего назначения.

## ⚙️ Технологии

- **Language:** [Go](https://go.dev/) (1.24+)
- **Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Auth:** JWT (JSON Web Tokens) & Bcrypt для хеширования паролей
- **Containerization:** Docker & Docker Compose