package vless

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// Config представляет собой структуру конфигурации для одного пользователя VLESS.
// Эта структура будет преобразована в JSON для sing-box.
type Config struct {
	Clients  []Client `json:"clients"`
	Dest     string   `json:"dest"`
	Flow     string   `json:"flow"`
	Protocol string   `json:"protocol"`
}

// Client представляет одного пользователя в конфигурации.
type Client struct {
	ID    string `json:"id"` // UUID пользователя
	Email string `json:"email"`
}

// GenerateVLESSUUID создает новый UUID, который используется как идентификатор пользователя в VLESS.
func GenerateVLESSUUID() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации UUID: %w", err)
	}
	return newUUID.String(), nil
}

// GenerateShortID генерирует короткий идентификатор (8 байт в hex) для XTLS-Reality.
func GenerateShortID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("ошибка при генерации short id: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

// GenerateKeyPair генерирует пару ключей x25519 для XTLS-Reality.
// Возвращает privateKey, publicKey, error
func GenerateKeyPair() (string, string, error) {
	// В Go нет нативной функции для генерации ключей x25519 в нужном формате.
	// Sing-box использует для этого свою внутреннюю библиотеку.
	// На данном этапе мы можем использовать заглушки или вызвать внешний скрипт.
	// Для MVP мы сгенерируем их с помощью `sing-box` установленного в системе.
	
	// ВНИМАНИЕ: Это временное решение для PoC. 
	// В проде нужно будет либо найти Go-библиотеку, либо использовать gRPC вызов к sing-box.
	// Здесь мы вернем статичные ключи-заглушки.
	
	// Пример вызова для генерации реальных ключей:
	// cmd := exec.Command("sing-box", "generate", "reality-keypair")
	
	return "dummy_private_key_replace_me", "dummy_public_key_replace_me", nil
}

// BuildVLESSLink создает строку подключения VLESS на основе предоставленных данных.
func BuildVLESSLink(uuid, publicKey, shortID, serverAddress, serverPort, serverName, sni string) string {
	// Формат ссылки VLESS-Reality:
	// vless://{uuid}@{server_addr}:{port}?encryption=none&flow=xtls-rprx-vision&security=reality&sni={sni}&fp=chrome&pbk={publicKey}&sid={shortId}&type=tcp...
	// Параметр `pbk` — это публичный ключ (Public Key).
	link := fmt.Sprintf("vless://%s@%s:%s?encryption=none&flow=xtls-rprx-vision&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp&headerType=none#%s",
		uuid,
		serverAddress,
		serverPort,
		sni,
		publicKey,
		shortID,
		serverName,
	)
	return link
}


// ToJSON преобразует структуру Config в JSON-строку.
func (c *Config) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге JSON: %w", err)
	}
	return string(data), nil
}
