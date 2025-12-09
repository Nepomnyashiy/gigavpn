package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// --- User Domain ---

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TelegramID   int64     `db:"telegram_id" json:"telegram_id"`
	Username     string    `db:"username" json:"username"`
	FullName     string    `db:"full_name" json:"full_name"`
	LanguageCode string    `db:"language_code" json:"language_code"`
	Balance      int64     `db:"balance" json:"balance"` // Храним в копейках
	IsBanned     bool      `db:"is_banned" json:"is_banned"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Transaction struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Amount    int64     `db:"amount" json:"amount"` // + Пополнение, - Списание
	Source    string    `db:"source" json:"source"` // "crypto", "card", "system"
	Status    string    `db:"status" json:"status"` // "pending", "success"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// --- Infrastructure Domain ---

type Server struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	IPAddress   string `db:"ip_address" json:"ip_address"` // Postgres inet -> string
	CountryCode string `db:"country_code" json:"country_code"`
	ApiURL      string `db:"api_url" json:"api_url"` // gRPC endpoint: "10.0.0.1:8080"

	// XrayConfig хранит специфичные настройки (порт, тип flow, сни)
	// Используем RawMessage, чтобы не парсить лишний раз, если нам нужно просто отдать JSON
	XrayConfig json.RawMessage `db:"xray_config" json:"xray_config"`

	CapacityLimit int  `db:"capacity_limit" json:"capacity_limit"`
	CurrentLoad   int  `db:"current_load" json:"current_load"`
	IsActive      bool `db:"is_active" json:"is_active"`
}

// --- Product Domain ---

type Plan struct {
	ID             int    `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Description    string `db:"description" json:"description"`
	Price          int64  `db:"price" json:"price"`
	DurationDays   int    `db:"duration_days" json:"duration_days"`
	TrafficLimitGB *int   `db:"traffic_limit_gb" json:"traffic_limit_gb"` // Pointer, т.к. может быть NULL (безлимит)
	MaxDevices     int    `db:"max_devices" json:"max_devices"`
}

type Subscription struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	PlanID    int       `db:"plan_id" json:"plan_id"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	AutoRenew bool      `db:"auto_renew" json:"auto_renew"`
}

// --- Technical Domain ---

type AccessKey struct {
	ID             uuid.UUID `db:"id" json:"id"` // Это VLESS UUID
	SubscriptionID uuid.UUID `db:"subscription_id" json:"subscription_id"`
	ServerID       int       `db:"server_id" json:"server_id"`

	KeyPrivate string `db:"key_private" json:"-"` // Не отдаем в JSON по умолчанию (безопасность)
	KeyPublic  string `db:"key_public" json:"key_public"`
	AccessLink string `db:"access_link" json:"access_link"`

	IsEnabled bool `db:"is_enabled" json:"is_enabled"`
}

type TrafficLog struct {
	ID            int64     `db:"id" json:"id"`
	AccessKeyID   uuid.UUID `db:"access_key_id" json:"access_key_id"`
	ServerID      int       `db:"server_id" json:"server_id"`
	UploadBytes   int64     `db:"upload_bytes" json:"upload_bytes"`
	DownloadBytes int64     `db:"download_bytes" json:"download_bytes"`
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`
}
