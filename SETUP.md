# Настройка проекта URL Shortener

## Чек-лист для запуска проекта

### 1. Требуемые версии ПО

| Компонент | Минимальная версия | Рекомендуемая версия | Проверка установки |
|-----------|-------------------|----------------------|-------------------|
| **Go** | 1.21 | 1.22+ | `go version` |
| **PostgreSQL** | 14 | 15+ | `psql --version` |
| **Docker** | 20.10 | 24+ | `docker --version` |
| **Docker Compose** | 2.0 | 2.20+ | `docker compose version` |
| **Git** | 2.30 | 2.40+ | `git --version` |

### 2. Установка Go-пакетов

Все зависимости уже указаны в `go.mod`. Для установки выполните:

```bash
# Скачать все зависимости
go mod download

# Или обновить зависимости
go mod tidy
```

**Основные зависимости:**
- `gorm.io/gorm` + `gorm.io/driver/postgres` - ORM для PostgreSQL
- `go.uber.org/zap` - структурированное логирование
- `github.com/spf13/viper` - управление конфигурацией
- `github.com/joho/godotenv` - загрузка .env файлов
- `github.com/golang-jwt/jwt/v5` - JWT аутентификация (для будущего использования)

### 3. Переменные окружения (.env)

Создайте файл `.env.local` на основе `.env.example`:

```bash
cp .env.example .env.local
```

**Обязательные переменные:**
```env
# Сервер
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_ENV=development

# База данных
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=urlshortener
DATABASE_PASSWORD=password
DATABASE_NAME=url_shortener
DATABASE_SSL_MODE=disable

# Аутентификация (опционально для MVP)
AUTH_JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

**Для продакшена измените:**
- `AUTH_JWT_SECRET` на случайную строку (минимум 32 символа)
- `DATABASE_PASSWORD` на надежный пароль
- `SERVER_ENV=production`

**Примечание:** Если вы используете Docker Compose для базы данных (порт 5433), установите `DATABASE_HOST=localhost` и `DATABASE_PORT=5433` в `.env.local`.

### 4. Запуск PostgreSQL локально

#### Вариант A: Docker Compose (рекомендуется)
```bash
# Запустить PostgreSQL
docker-compose up -d postgres

# Проверить статус
docker-compose ps

# Остановить
docker-compose down
```

#### Вариант B: Docker run
```bash
docker run -d \
  --name url-shortener-postgres \
  -e POSTGRES_USER=urlshortener \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=url_shortener \
  -p 5432:5432 \
  postgres:15-alpine
```

#### Вариант C: Локальная установка PostgreSQL
1. Установите PostgreSQL с официального сайта
2. Создайте базу данных:
```sql
CREATE DATABASE url_shortener;
CREATE USER urlshortener WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE url_shortener TO urlshortener;
```

### 5. Расширения для VSCode

**Обязательные:**
- [Go](https://marketplace.visualstudio.com/items?itemName=golang.go) - поддержка Go
- [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) - работа с Docker
- [GitLens](https://marketplace.visualstudio.com/items?itemName=eamodio.gitlens) - улучшенный Git

**Рекомендуемые:**
- [Go Test Explorer](https://marketplace.visualstudio.com/items?itemName=premparihar.gotestexplorer) - навигация по тестам
- [YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) - редактирование YAML файлов
- [Markdown All in One](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one) - работа с Markdown
- [Error Lens](https://marketplace.visualstudio.com/items?itemName=usernamehw.errorlens) - подсветка ошибок

### 6. Настройка Go в VSCode

1. Откройте командную палитру (Ctrl+Shift+P)
2. Выполните `Go: Install/Update Tools`
3. Выберите все инструменты или:
   - `gopls` - language server
   - `dlv` - debugger
   - `staticcheck` - статический анализатор
   - `golangci-lint` - линтер

### 7. Первый запуск проекта

Есть два основных способа запуска:

#### Способ A: Локальный запуск приложения с Docker-контейнерами БД

```bash
# 1. Установить зависимости
go mod download

# 2. Настроить окружение
cp .env.example .env.local
# Отредактируйте .env.local:
# DATABASE_HOST=localhost
# DATABASE_PORT=5433   # порт, проброшенный из контейнера

# 3. Запустить базу данных
docker-compose up -d postgres

# 4. Запустить миграции
make migrate

# 5. Запустить приложение локально
go run ./cmd/server

# Или с hot reload (требуется air)
make dev
```

#### Способ B: Полный запуск через Docker Compose (приложение тоже в контейнере)

```bash
# Запустить все сервисы (PostgreSQL, приложение)
docker-compose up -d

# Проверить логи приложения
docker-compose logs -f app

# Остановить
docker-compose down
```

#### Способ C: Локальный запуск с локальной БД (без Docker)

1. Установите PostgreSQL локально, создайте базу `url_shortener` и пользователя `urlshortener`.
2. В `.env.local` укажите `DATABASE_HOST=localhost`, `DATABASE_PORT=5432`.
3. Запустите миграции: `make migrate`.
4. Запустите приложение: `go run ./cmd/server`.

**Примечание:** По умолчанию health-эндпоинт доступен по адресу http://localhost:8080/health.

### 8. Проверка работоспособности

1. **Health check:** http://localhost:8080/health (должен вернуть "OK")
2. **Логи:** Проверьте вывод в консоли
3. **База данных:** Подключитесь к PostgreSQL:
```bash
psql -h localhost -U urlshortener -d url_shortener
```

4. **Тестирование API:**
```bash
# Создать короткую ссылку
curl -X POST http://localhost:8080/api/v1/links \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'

# Получить статистику
curl http://localhost:8080/api/v1/links/{code}/stats

# Перейти по короткой ссылке (в браузере)
open http://localhost:8080/{code}
```

### 9. Полезные команды

```bash
# Сборка
make build

# Запуск тестов
make test

# Линтинг
make lint

# Запуск в Docker
make docker-up

# Остановка Docker
make docker-down

# Очистка
make clean

# Запуск миграций
make migrate
```

### 10. Устранение неполадок

#### Ошибка: "connection refused" к PostgreSQL
- Проверьте, что PostgreSQL запущен: `docker-compose ps`
- Проверьте порт: `netstat -an | grep 5432` (или 5433)
- Проверьте credentials в `.env.local`

#### Ошибка: "missing go.sum entry"
```bash
go mod tidy
go mod download
```

#### Ошибка: "permission denied" при запуске Docker
- Добавьте пользователя в группу docker: `sudo usermod -aG docker $USER`
- Перезайдите в систему

#### Ошибка: "go: cannot find main module"
- Убедитесь, что вы в корневой директории проекта
- Проверьте наличие `go.mod`

#### Health-эндпоинт не отвечает
- Убедитесь, что приложение запущено: `curl -v http://localhost:8080/health`
- Проверьте логи приложения: `docker-compose logs app`
- Убедитесь, что сервер слушает на 0.0.0.0:8080

### 11. Дальнейшие шаги

После успешного запуска:
1. Изучите структуру проекта в `plans/`
2. Протестируйте основные эндпоинты API
3. Добавьте аутентификацию (если требуется)
4. Добавьте тесты
5. Настройте мониторинг и деплой

### 12. Дополнительные ресурсы

- [Документация Go](https://golang.org/doc/)
- [Документация GORM](https://gorm.io/docs/)
- [Документация PostgreSQL](https://www.postgresql.org/docs/)
- [Документация Zap](https://pkg.go.dev/go.uber.org/zap)

---

**Готово!** Проект должен запускаться на http://localhost:8080. Для начала разработки следуйте roadmap в `plans/development-roadmap.md`.