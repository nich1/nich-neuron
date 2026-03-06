from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    OLLAMA_URL: str = "http://localhost:11434"
    OPEN_WEBUI_URL: str = "http://localhost:3000"
    N8N_URL: str = "http://localhost:5678"
    API_TOKEN: str = ""


settings = Settings()