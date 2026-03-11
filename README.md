# nich-neuron

Self-hosted AI platform running on a single home server. Everything runs in Docker Compose.

## Prerequisites

- **Docker Engine** 24+ with the **Compose v2** plugin
- **4 GB RAM** minimum (8+ recommended — LLMs are memory-hungry)
- **NVIDIA GPU** optional — detected automatically by `setup.sh`

## Quick Start

```bash
git clone https://github.com/nich1/nich-neuron.git
cd nich-neuron
./setup.sh
docker compose up -d
```

`setup.sh` checks prerequisites, generates secrets, detects your GPU, and creates `.env`. After that, `docker compose up -d` starts the full stack.

To set up manually instead, copy `sample.env` to `.env` and replace every `CHANGE_ME` placeholder with a random value.

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
| Vaultwarden    | 8443  | vaultwarden/server                                          |
| Uptime Kuma    | 3001  | louislam/uptime-kuma                                        |
| Dashboard      | 9000  | local build (nginx)                                         |

## Models

Edit `models.txt` to add or remove Ollama models. They are pulled automatically on first startup by the `ollama-init` container.

## GPU Passthrough

`setup.sh` detects NVIDIA GPUs and offers to enable passthrough. This creates a `docker-compose.override.yml` with the NVIDIA runtime config. To enable manually, create that file with:

```yaml
services:
  ollama:
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]
```

## Checking Health

Every service has a Docker healthcheck. After starting the stack:

```bash
docker compose ps
```

The `STATUS` column shows `healthy`, `starting`, or `unhealthy` for each container.
