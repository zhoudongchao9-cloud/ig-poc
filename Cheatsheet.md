# Inference Gateway – Cheatsheet

## Quick Commands

| Action | Command |
|--------|---------|
| Run with OpenAI | `docker run -p 8080:8080 -e OPENAI_API_KEY=xxx ghcr.io/inference-gateway/inference-gateway:latest` |
| List models | `curl http://localhost:8080/v1/models` |
| Chat (OpenAI) | `curl -X POST http://localhost:8080/v1/chat/completions -H "Content-Type: application/json" -d '{"model":"openai/gpt-5","messages":[{"role":"user","content":"Hi"}]}'` |
| Chat (Anthropic) | same, with `"model":"anthropic/claude-opus-4-8"` |
| Enable vision | set `ENABLE_VISION=true` |

## Common Environment Varieables

| Variable | Purpose |
|----------|---------|
| `OPENAI_API_KEY` | OpenAI API key |
| `ANTHROPIC_API_KEY` | Anthropic key |
| `GROQ_API_KEY` | Groq key |
| `AUTH_ENABLE` | Enable OIDC (default false) |
| `PORT` | Server port (default 8080) |
| `LOG_LEVEL` | debug / info / warn / error |

## Health Check

```bash
curl http://localhost:8080/health
