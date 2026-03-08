# nich-neuron

Home server / lab setup.

Monorepo that includes the code for a custom agent platform as well as a docker file for
wiring everything up. Comes with n8n, Open WebUI, Agent Platform API (in development), and Ollama.
Agent Framework Repo:

## Quick start

1. Copy env and set secrets:
   ```bash
   cp sample.env .env
   # Edit .env: set N8N_PASSWORD, WEBUI_SECRET_KEY, optional API_TOKEN
   ```
2. Start everything:
   ```bash
   docker compose up -d
   ```
3. Open:
   - **Open WebUI (chat):** http://localhost:3000
   - **n8n:** http://localhost:5678
   - **Agent platform API:** http://localhost:8000 (docs at `/docs`)
   - **Ollama:** http://localhost:11434
   - **Qdrant (vector DB):** http://localhost:6333 (for RAG/semantic search when you add it)

4. Models: `ollama-init` pulls models listed in `models.txt` on first startup. Edit that file to add or remove models (one per line, `#` for comments).
