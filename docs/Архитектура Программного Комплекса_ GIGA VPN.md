Архитектура Программного Комплекса: GIGA VPN
Тип документа: System Architecture Design Document (SADD)
Стек: Go (Backend) + Python (Frontend) + PostgreSQL + Sing-box
Архитектурный стиль: Microservices / Service-Oriented Architecture (SOA)

1. Концептуальная схема (High-Level Architecture)
Система разделена на три логических уровня, чтобы обеспечить изоляцию ответственности и высокую надежность.
1.1. Уровень Интерфейса (Frontend Layer)
    • Компонент: VPN-Bot
    • Технология: Python 3.11 + Aiogram 3.
    • Ответственность:
        ◦ Взаимодействие с пользователем (UI/UX).
        ◦ Обработка команд (/start, /pay, /help).
        ◦ Отображение QR-кодов и инструкций.
        ◦ Ограничение: Бот не имеет прямого доступа к базе данных. Все данные он получает через API Бэкенда.
1.2. Уровень Управления (Control Plane / Core)
    • Компонент: VPN-Orchestrator
    • Технология: Go 1.22 + Gin (HTTP) + gRPC.
    • Ответственность:
        ◦ User Management: Регистрация, баланс, подписки.
        ◦ Node Management: Выбор сервера для размещения пользователя (Load Balancing).
        ◦ Provisioning: Отправка команд на удаленные VPN-ноды (создать/удалить ключ).
        ◦ Billing: Обработка транзакций.
1.3. Уровень Данных (Data Plane / Infrastructure)
    • Компонент: VPN-Node
    • Технология: Linux + Docker + Sing-box/Xray.
    • Ответственность:
        ◦ Непосредственная маршрутизация и шифрование трафика.
        ◦ Сбор метрик использования (Traffic Accounting).

2. Структура Проекта (Repository Layout)
Мы будем использовать структуру Monorepo (один репозиторий для удобства на старте), разделенный на сервисы. Это стандарт индустрии.
Plaintext
/vpn-platform-root
├── /backend-go               # <-- СЕРДЦЕ СИСТЕМЫ (Control Plane)
│   ├── /cmd
│   │   └── /server           # Точка входа (main.go)
│   ├── /internal             # Приватный код приложения
│   │   ├── /config           # Чтение конфигов (env, yaml)
│   │   ├── /domain           # Бизнес-сущности (User, Server, Subscription)
│   │   ├── /service          # Бизнес-логика (UserService, NodeService)
│   │   ├── /repository       # Работа с БД (PostgreSQL / SQLc)
│   │   └── /transport        # API Хендлеры
│   │       ├── /http         # REST API для Бота
│   │       └── /grpc         # gRPC клиент для управления Нодами
│   ├── /pkg                  # Библиотеки, доступные извне (утилиты)
│   ├── /migrations           # SQL файлы миграции БД
│   ├── go.mod                # Зависимости Go
│   └── Dockerfile            # Сборка Go-сервиса
│
├── /bot-python               # <-- ИНТЕРФЕЙС (Frontend)
│   ├── /handlers             # Обработчики команд телеграма
│   ├── /keyboards            # Кнопки и меню
│   ├── /services             # Клиент API к backend-go
│   ├── bot.py                # Запуск бота
│   └── requirements.txt      # Зависимости Python
│
├── /infrastructure           # <-- DevOps (IaC)
│   ├── /terraform            # Скрипты для создания серверов (Aeza/DO)
│   ├── /ansible              # Скрипты настройки (установка Xray, Docker)
│   └── /vpn-node-config      # Шаблоны конфигов Sing-box/Xray
│
├── /proto                    # <-- Контракты API
│   └── vpn_service.proto     # Описание gRPC взаимодействия
│
└── docker-compose.yml        # Локальный запуск (БД + Бэкенд + Бот)


3. Процессы обработки данных (Data Flow)
Сценарий 1: Создание нового VPN-ключа (Provisioning)
Этот процесс показывает, как Бот, Бэкенд и Сервер взаимодействуют друг с другом.
    1. Bot: Пользователь нажимает "Получить доступ".
        ◦ Бот отправляет POST /api/subs/create на Go Backend.
    2. Go Backend (Service Layer):
        ◦ Проверяет баланс пользователя в БД.
        ◦ Алгоритм NodeSelector выбирает наименее нагруженный сервер (например, node-nl-01).
        ◦ Генерирует пару ключей (Public/Private) и UUID.
    3. Go Backend (Infrastructure Layer):
        ◦ Устанавливает gRPC/SSH соединение с node-nl-01.
        ◦ Вызывает метод AddClient(uuid, private_key).
    4. VPN Node:
        ◦ Применяет изменения в памяти (hot reload).
        ◦ Возвращает "OK".
    5. Go Backend:
        ◦ Сохраняет запись в PostgreSQL: Subscription {User: 1, Node: node-nl-01, Key: ...}.
        ◦ Генерирует ссылку vless://....
        ◦ Отдает JSON ответ Боту.
    6. Bot: Генерирует QR-код из ссылки и шлет юзеру.

4. Стек технологий и Инструменты
Backend (Go)
    • Framework: Gin (самый быстрый HTTP веб-фреймворк) или Chi.
    • DB Driver: pgx (высокопроизводительный драйвер PostgreSQL).
    • SSH Client: crypto/ssh (нативная библиотека Go для управления серверами).
    • Logging: slog (стандартный логгер Go 1.21+).
Database
    • PostgreSQL 15+: Основное хранилище.
    • Redis 7: Кеширование сессий и кратковременные данные.
DevOps
    • Docker: Контейнеризация сервисов.
    • Make: Автоматизация команд (make run, make build).

5. План разработки (Implementation Path)
Мы будем двигаться от ядра к периферии.
    1. Этап 1: "MVP Core" (Бэкенд + 1 Сервер).
        ◦ Создаем структуру папок Go.
        ◦ Реализуем модуль подключения к серверу по SSH/API.
        ◦ Реализуем генерацию VLESS-ссылок.
        ◦ Результат: Консольная программа на Go, которая создает рабочий ключ на сервере Aeza.
    2. Этап 2: "Persistence" (База данных).
        ◦ Поднимаем PostgreSQL.
        ◦ Реализуем сохранение пользователей.
        ◦ Поднимаем HTTP API.
    3. Этап 3: "Interface" (Бот).
        ◦ Пишем простого бота на Python, который стучится в наш API.

Вердикт Архитектора
Предложенная структура monorepo с четким разделением backend-go и bot-python позволяет:
    1. Менять логику VPN (например, сменить протокол VLESS на Trojan), не переписывая Бота.
    2. Менять дизайн Бота, не трогая ядро системы.
    3. Легко масштабировать команду (один пишет на Go, другой на Python).