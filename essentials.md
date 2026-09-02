<!-- This file is temporary -->

# Inference Gateway

Inference Gateway is a lightweight, cloud‑native proxy that unifies multiple LLM providers behind a single OpenAI‑compatible API. It simplifies access to various language models, supports streaming, tool calling, and MCP, and is designed for production use with observability, authentication, and Kubernetes support.

---

## Key Features

- **Unified API** – One OpenAI‑compatible endpoint for OpenAI, Anthropic, Groq, Cohere, Ollama, llama.cpp, Google, Mistral, and many more.
- **Tool & Function Calling** – Built‑in support for tool use across supported providers.
- **MCP Integration** – Automatic discovery and exposure of MCP server tools to LLMs.
- **Streaming & Vision** – Real‑time token streaming and multi‑modal (image+text) support.
- **Observability** – Prometheus metrics (GenAI semantic conventions) and OTLP export.
- **Enterprise Ready** – OIDC authentication, authorisation, configurable timeouts, and TLS.
- **Lightweight** – ~10.8 MB binary, minimal resource footprint.
- **Self‑hosted** – Zero data collection, Apache 2.0 license.

---

## Quick Start

### Docker

```bash
docker pull ghcr.io/inference-gateway/inference-gateway:latest
docker run --rm -it -p 8080:8080 -e OPENAI_API_KEY=your_key_here \
  ghcr.io/inference-gateway/inference-gateway:latest
```

### Docker Compose

```yaml
services:
  inference-gateway:
    image: ghcr.io/inference-gateway/inference-gateway:latest
    environment:
      - OPENAI_API_KEY=your-api-key
    ports:
      - '8080:8080'
```

### Kubernetes

See the [Operator](https://docs.inference-gateway.com/operator/) for declarative management. Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inference-gateway
spec:
  replicas: 2
  selector:
    matchLabels:
      app: inference-gateway
  template:
    metadata:
      labels:
        app: inference-gateway
    spec:
      containers:
        - name: inference-gateway
          image: ghcr.io/inference-gateway/inference-gateway:latest
          ports:
            - containerPort: 8080
          env:
            - name: OPENAI_API_KEY
              valueFrom:
                secretKeyRef:
                  name: llm-secrets
                  key: openai-api-key
```

---

## Supported Providers

- OpenAI, Anthropic, Groq, Cohere
- Ollama, llama.cpp, Cloudflare
- DeepSeek, Google, Mistral
- MiniMax, Moonshot, Nvidia, Z.ai

---

## Basic API Usage

The gateway runs at `http://localhost:8080` by default.

### List available models

```bash
curl http://localhost:8080/v1/models
```

### Chat completion (OpenAI)

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai/gpt-5",
    "messages": [
      { "role": "system", "content": "You are a helpful assistant." },
      { "role": "user", "content": "Explain how Inference Gateway works." }
    ]
  }'
```

### Using Anthropic

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "anthropic/claude-opus-4-8",
    "messages": [
      { "role": "system", "content": "You are a helpful assistant." },
      { "role": "user", "content": "Compare different LLM providers." }
    ]
  }'
```

### Vision support (enable with `ENABLE_VISION=true`)

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "anthropic/claude-opus-4-8",
    "messages": [
      {
        "role": "user",
        "content": [
          { "type": "text", "text": "What is in this image?" },
          { "type": "image_url", "image_url": { "url": "https://example.com/image.jpg" } }
        ]
      }
    ]
  }'
```

---

## Deployment Options

- **Binary** – Download and run directly.
- **Docker / Docker Compose** – Recommended for containerised environments.
- **Kubernetes** – Use the official Operator for production‑grade deployments.

---

## Additional Resources

- **GitHub** – [github.com/inference-gateway/inference-gateway](https://github.com/inference-gateway/inference-gateway)
- **Documentation** – [docs.inference-gateway.com](https://docs.inference-gateway.com)
- **CLI Tool** – [github.com/inference-gateway/cli](https://github.com/inference-gateway/cli)
- **Operator** – [github.com/inference-gateway/operator](https://github.com/inference-gateway/operator)

---

## License

Apache 2.0
```

Feel free to copy and modify this Markdown as needed.
