# Dockerfile для Go бэкенда

# Этап сборки
FROM golang:1.22-alpine AS builder

# Установка зависимостей, необходимых для сборки (если есть)
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
# CGO_ENABLED=0 для статической компиляции, чтобы бинарник не зависел от системных библиотек
# -o server - имя выходного файла
# ./cmd/server - путь к main пакету
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Этап выполнения
FROM alpine:latest

# Устанавливаем сертификаты CA для HTTPS-соединений (например, для базы данных)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/server .

# Копируем файлы миграций
COPY --from=builder /app/migrations ./migrations

# Устанавливаем переменную окружения для буферизации вывода Go
ENV GODEBUG=asyncpreemptoff=1

# Открываем порт, который слушает Go-приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./server"]
