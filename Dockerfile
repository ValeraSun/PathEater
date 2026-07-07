# Базовый образ с Go
FROM golang:1.23-alpine

# Устанавливаем полезные утилиты (опционально)
RUN apk add --no-cache git curl

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования слоев
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
#RUN go build -o /app/main ./cmd/main.go

# Открываем порт (замените на ваш, если нужно)
EXPOSE 8080

# Запускаем приложение
#CMD ["/app/main"]