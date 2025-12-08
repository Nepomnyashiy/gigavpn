import httpx
import asyncio

async def check_backend():
    print("Запуск проверки бэкенда с помощью httpx...")
    # Используем те же параметры, что и в основном приложении
    backend_url = "http://gigavpn_backend:8080"
    
    try:
        async with httpx.AsyncClient(trust_env=False) as client:
            print(f"Отправка GET-запроса на {backend_url}/ping")
            response = await client.get(f"{backend_url}/ping")
            print(f"Ответ получен. Статус-код: {response.status_code}")
            response.raise_for_status()
            print("Ответ успешный. Тело ответа:")
            print(response.text)
    except Exception as e:
        print(f"!!! Произошла ошибка httpx:")
        print(e)

if __name__ == "__main__":
    asyncio.run(check_backend())
