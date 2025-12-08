package main

import (
	"fmt"
	"log"

	"gigavpn/backend-go/internal/vless"
)

func main() {
	fmt.Println("Запуск VPN-Orchestrator (Control Plane)...")
	fmt.Println("Этап 1: Генерация тестового ключа VLESS...")

	// 1. Генерируем UUID для пользователя
	vlessUUID, err := vless.GenerateVLESSUUID()
	if err != nil {
		log.Fatalf("Не удалось сгенерировать VLESS UUID: %v", err)
	}
	fmt.Printf(" - Сгенерирован VLESS UUID: %s\n", vlessUUID)

	// 2. Генерируем пару ключей Reality (пока используются заглушки)
	// Приватный ключ сохраняется в БД, а публичный используется в ссылке.
	_, publicKey, err := vless.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Не удалось сгенерировать ключи: %v", err)
	}
	fmt.Printf(" - Сгенерирована пара ключей (Public): %s\n", publicKey)

	// 3. Генерируем ShortID для Reality
	shortID, err := vless.GenerateShortID()
	if err != nil {
		log.Fatalf("Не удалось сгенерировать ShortID: %v", err)
	}
	fmt.Printf(" - Сгенерирован ShortID: %s\n", shortID)

	// 4. Собираем финальную ссылку подключения
	// Эти данные обычно берутся из БД или конфига
	const serverAddress = "192.168.1.1"
	const serverPort = "443"
	const serverName = "MyVPN-NL"
	const sni = "www.microsoft.com"

	link := vless.BuildVLESSLink(vlessUUID, publicKey, shortID, serverAddress, serverPort, serverName, sni)

	fmt.Println("\n--- РЕЗУЛЬТАТ ---")
	fmt.Println("Полная ссылка для подключения (VLESS):")
	fmt.Printf("%s\n", link)
	fmt.Println("\nЭтап 1 (MVP Core) успешно продемонстрирован.")
	fmt.Println("Следующие шаги: подключение к БД, реализация API и автоматизация через SSH/gRPC.")
}
