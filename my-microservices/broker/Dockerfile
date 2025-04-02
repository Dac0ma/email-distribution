# Этап сборки
FROM golang:1.20-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Компилируем приложение
RUN go build -o broker .

# Этап выполнения
FROM alpine:latest

# Установка необходимых библиотек (если нужно)
RUN apk --no-cache add ca-certificates

# Установка рабочей директории
WORKDIR /root/

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/broker .

# Команда для запуска приложения
CMD ["./broker"]