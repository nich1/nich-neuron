"""Nich Neuron Home Server Endpoint"""
from fastapi import FastAPI

from .settings import settings

app = FastAPI(title="Agent Platform", version="0.1.0")


@app.get("/")
async def root():
    return {"service": "nich-neuron", "status": "ok", "docs": "/docs"}


@app.get("/health")
async def health():
    return {"status": "ok", "ollama": settings.OLLAMA_URL, "n8n": settings.N8N_URL, "open_webui": settings.OPEN_WEBUI_URL}
