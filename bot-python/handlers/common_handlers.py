from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command
import asyncio

from services.api_client import api_client

# Создаем новый роутер для общих команд
router = Router()

@router.message(Command("start"))
async def handle_start(message: Message):
    """
    Обработчик команды /start.
    Приветствует пользователя и проверяет состояние бэкенда.
    """
    # Добавляем небольшую задержку на случай проблем с race condition в сети Docker
    await asyncio.sleep(1)
    
    # Проверяем доступность бэкенда
    is_backend_ok = await api_client.ping_server()
    
    if is_backend_ok:
        backend_status = "✅ Онлайн"
    else:
        backend_status = "❌ Оффлайн"
        
    welcome_text = (
        f"👋 **Добро пожаловать, {message.from_user.full_name}!**\n\n"
        f"Я — бот для управления GigaVPN.\n\n"
        f"Статус системы:\n"
        f" - Бэкенд: {backend_status}\n\n"
        f"Используйте /help для просмотра доступных команд."
    )
    
    await message.answer(welcome_text, parse_mode="Markdown")
