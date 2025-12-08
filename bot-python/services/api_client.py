import httpx
import logging

# URL нашего Go-бэкенда.
# В реальном приложении это значение лучше выносить в переменные окружения.
BACKEND_URL = "http://localhost:8080"

class APIClient:
    """
    HTTP-клиент для взаимодействия с API Go-бэкенда.
    """
    def __init__(self, base_url: str):
        self.base_url = base_url
        # Создаем асинхронный клиент, который будет использоваться для всех запросов.
        # trust_env=False отключает использование системных прокси и переменных окружения.
        self.client = httpx.AsyncClient(base_url=self.base_url, timeout=10.0, trust_env=False)

    async def ping_server(self) -> bool:
        """
        Проверяет доступность бэкенд-сервера.
        :return: True, если сервер доступен и отвечает "pong", иначе False.
        """
        try:
            response = await self.client.get("/ping")
            
            # Проверяем, что запрос успешен (статус код 2xx)
            response.raise_for_status()
            
            # Проверяем, что в ответе есть ожидаемое сообщение
            if response.json().get("message") == "pong":
                logging.info("Сервер бэкенда успешно ответил на ping.")
                return True
            else:
                logging.warning("Сервер бэкенда ответил, но тело ответа некорректно.")
                return False

        except httpx.RequestError as e:
            logging.error(f"Ошибка при запросе к бэкенду: {e}")
            return False
        except Exception as e:
            logging.error(f"Непредвиденная ошибка при проверке сервера: {e}")
            return False

# Создаем глобальный экземпляр клиента, который будет использоваться во всем приложении бота.
api_client = APIClient(base_url=BACKEND_URL)
