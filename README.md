# My Blog Backend

Backend-сервис для платформы блога, предоставляющий API для управления пользователями, статьями и категориями. Проект построен с использованием чистой архитектуры и современных подходов к разработке на Go.

## Быстрый старт

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

## Структура проекта

- `cmd/app`: Точка входа в приложение.
- `db/migrations`: SQL файлы для миграции базы данных.
- `internal/config`: Логика загрузки конфигурации.
- `internal/delivery/v1`: Обработка HTTP запросов (Handlers & Middleware).
- `internal/repository`: Работа с базой данных (GORM).
- `internal/domain`: Доменные сущности.
- `internal/usecase`: Бизнес-логика приложения.
- `pkg/`: Вспомогательные утилиты и пакеты общего назначения.

## Технологии

- **Language:** [Go](https://go.dev/) (1.24+)
- **Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Auth:** JWT (JSON Web Tokens) & Bcrypt для хеширования паролей
- **Containerization:** Docker & Docker Compose

---