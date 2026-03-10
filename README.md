# nich-neuron

Home server / lab setup.

## Services

| Service        | Port  | Image                                                       |
| -------------- | ----- | ----------------------------------------------------------- |
| Ollama         | 11434 | ollama/ollama                                               |
| Open WebUI     | 3000  | ghcr.io/open-webui/open-webui                               |
| n8n            | 5678  | n8nio/n8n                                                   |
| Agent Platform | 8000  | [neuron-orion](https://github.com/nich1/neuron-orion)       |
| Auth Server    | 8081  | [neuron-cerberus](https://github.com/nich1/neuron-cerberus) |
| Qdrant         | 6333  | qdrant/qdrant                                               |
| Seq            | 5380  | datalust/seq                                                |
| Ntfy           | 8090  | binwiederhier/ntfy                                          |
| Dashboard      | 9000  | local build (nginx)                                         |

## Quick Start

1. Copy env and set secrets:
   ```bash
   cp sample.env .env
   ```
2. Start everything:
   ```bash
   docker compose up -d
   ```
3. Edit `models.txt` to add or remove Ollama models (pulled automatically on first startup).
