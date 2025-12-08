import asyncio
import logging
import os

from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.types import BotCommand

# Импортируем наш обработчик
from handlers import common_handlers

# Настройка логирования
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

async def set_commands(bot: Bot):
    """
    Устанавливает команды, которые будут видны в меню Telegram.
    """
    commands = [
        BotCommand(command="/start", description="Запустить бота"),
        BotCommand(command="/help", description="Помощь"),
    ]
    await bot.set_my_commands(commands)

async def main():
    """
    Основная функция для запуска бота.
    """
    logging.info("Старт функции main()")
    
    # Получаем токен бота из переменной окружения
    bot_token = os.getenv("BOT_TOKEN")
    if not bot_token:
        logging.error("Ошибка: не найден токен бота. Установите переменную окружения BOT_TOKEN.")
        return
    else:
        logging.info(f"Токен получен. Первые 8 символов: {bot_token[:8]}...")

    # Инициализация бота и диспетчера с использованием нового синтаксиса для parse_mode
    bot = Bot(token=bot_token, default=DefaultBotProperties(parse_mode=ParseMode.HTML))
    dp = Dispatcher()

    # Регистрация роутеров (обработчиков команд)
    dp.include_router(common_handlers.router)

    # Установка команд меню
    await set_commands(bot)

    logging.info("Запуск бота...")
    # Удаляем вебхук, если он был установлен ранее
    await bot.delete_webhook(drop_pending_updates=True)
    # Запускаем поллинг
    await dp.start_polling(bot)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except (KeyboardInterrupt, SystemExit):
        logging.info("Бот остановлен.")