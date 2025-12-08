import aiohttp
import logging
import os

# URL нашего Go-бэкенда.
# Получаем его из переменной окружения, как и положено в Docker.
BACKEND_URL = os.getenv("BACKEND_URL", "http://localhost:8080")

class APIClient:
    """
    HTTP-клиент для взаимодействия с API Go-бэкенда, использующий aiohttp.
    """
    def __init__(self, base_url: str):
        self.base_url = base_url
        # Сессия создается по требованию, а не в конструкторе.
        self._session = None

    async def get_session(self) -> aiohttp.ClientSession:
        """
        Создает или возвращает существующую сессию aiohttp.
        """
        if self._session is None or self._session.closed:
            # trust_env=False - важно, чтобы снова не подхватить системные прокси
            self._session = aiohttp.ClientSession(trust_env=False)
        return self._session

    async def close_session(self):
        """
        Закрывает сессию, если она была создана.
        """
        if self._session and not self._session.closed:
            await self._session.close()

    async def ping_server(self) -> bool:
        """
        Проверяет доступность бэкенд-сервера.
        :return: True, если сервер доступен и отвечает "pong", иначе False.
        """
        try:
            session = await self.get_session()
            url = f"{self.base_url}/ping"
            
            async with session.get(url, timeout=5) as response:
                logging.info(f"Запрос на {url}, получен статус: {response.status}")
                response.raise_for_status()
                
                data = await response.json()
                if data.get("message") == "pong":
                    logging.info("Сервер бэкенда успешно ответил на ping.")
                    return True
                else:
                    logging.warning("Сервер бэкенда ответил, но тело ответа некорректно.")
                    return False

        except aiohttp.ClientError as e:
            logging.error(f"Ошибка клиента aiohttp при запросе к бэкенду: {e}")
            return False
        except Exception as e:
            logging.error(f"Непредвиденная ошибка при проверке сервера: {e}")
            return False

# Создаем глобальный экземпляр клиента
api_client = APIClient(base_url=BACKEND_URL)
