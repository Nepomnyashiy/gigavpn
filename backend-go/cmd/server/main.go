package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gigavpn/backend-go/internal/repository"
	"gigavpn/backend-go/internal/transport/http"
)

func main() {
	fmt.Println("Запуск VPN-Orchestrator (Control Plane)...")

	// --- Конфигурация подключения к БД ---
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Переменная окружения DATABASE_URL не установлена.")
	}

	// --- Инициализация пула соединений с БД ---
	dbPool, err := repository.NewPostgresDB(context.Background(), databaseURL) // Теперь NewPostgresDB принимает URL
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Close()
	fmt.Println("✅ Успешное подключение к базе данных!")

	// --- Инициализация HTTP-сервера ---
	// В будущем мы передадим dbPool в NewHandler, чтобы обработчики могли работать с БД.
	handler := http.NewHandler()
	router := handler.InitRoutes()

	// --- Запуск сервера ---
	serverPort := "8080"
	log.Printf("Запуск HTTP-сервера на порту %s...", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %v", err)
	}
}
