# Dockerfile для Telegram-бота

# Используем официальный образ Python как базовый
FROM python:3.12-slim-bookworm

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Устанавливаем curl для отладки
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Копируем файл зависимостей и устанавливаем их
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копируем остальной код приложения
COPY . .

# Устанавливаем переменную окружения для буферизации вывода Python
# Это важно, чтобы логи выводились сразу, а не буферизовались
ENV PYTHONUNBUFFERED=1

# Открываем порт, если бот будет использовать webhook (пока неактуально, но на будущее)
# EXPOSE 8080 

# Запускаем бота
CMD ["python3", "bot.py"]
