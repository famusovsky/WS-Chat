### Выполнил Степанов Алексей Александрович

## Запуск:

Запуск dev-среды с помощью docker-compose:

```bash
docker-compose up
```

Запуск с помощью go run:

```bash
# Среда, в которой происходит запуск, должна иметь переменные окружения:
# DB_HOST
# DB_PORT
# DB_USER
# DB_PASSWORD
# DB_NAME
# API_PORT
go run ./cmd/server/main.go
# Флаги:
# -override_tables=true - запуск с автоматическим созданием таблиц в БД
# -addr=:8080 - выбор порта, с которым будет работать сервер
```

## PostgreSQL Query для создания таблиц в БД вручную:
```sql
CREATE TABLE messages (
	id SERIAL UNIQUE,
	text TEXT,
	sender TEXT,
	creation TIMESTAMP
);
```

