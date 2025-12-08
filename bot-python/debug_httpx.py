import httpx
import inspect

print(f"httpx version: {httpx.__version__}")
print(f"httpx file: {httpx.__file__}")

# Дополнительно проверим сигнатуру конструктора
sig = inspect.signature(httpx.AsyncClient.__init__)
print(f"httpx.AsyncClient.__init__ signature: {sig}")
