import asyncio
import logging
import os

# Настройка логирования
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

async def main():
    """
    Основная функция для запуска бота.
    Упрощена для диагностики проблемы с бесшумным падением.
    """
    logging.info("Старт функции main() - Упрощенная версия для отладки.")
    
    # Получаем токен бота из переменной окружения
    bot_token = os.getenv("BOT_TOKEN")
    if not bot_token:
        logging.error("Ошибка: не найден токен бота. Установите переменную окружения BOT_TOKEN.")
        return
    else:
        logging.info(f"Токен получен. Первые 8 символов: {bot_token[:8]}...")

    logging.info("Бот запущен. Ожидание 10 секунд перед выходом...")
    await asyncio.sleep(10) # Просто ждем 10 секунд
    logging.info("Бот завершил работу после 10 секунд.")


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except (KeyboardInterrupt, SystemExit):
        logging.info("Бот остановлен вручную.")
