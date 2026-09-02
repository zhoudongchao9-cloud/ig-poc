
# Inference Gateway – Advanced Configuration & Features

## Environment Variables (Key Ones)

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging verbosity (`debug`,`info`,`warn`,`error`) | `info` |
| `AUTH_ENABLE` | Enable OIDC authentication | `false` |
| `OIDC_ISSUER` | OIDC issuer URL | – |
| `OIDC_AUDIENCE` | Expected audience for JWT | – |
| `ENABLE_VISION` | Enable image input support | `false` |
| `REQUEST_TIMEOUT` | Upstream request timeout (seconds) | `30` |
| `TLS_ENABLE` | Enable HTTPS | `false` |
| `TLS_CERT_FILE`, `TLS_KEY_FILE` | Certificate and key paths | – |

## Authentication & Authorisation

- **OIDC** – Set `AUTH_ENABLE=true` and configure `OIDC_ISSUER` and `OIDC_AUDIENCE`.
- All API calls must include `Authorization: Bearer <JWT>` when enabled.
- JWT tokens are validated against the IdP; claims can be used for fine‑grained access control.

## Obsevrability

- **Prometheus metrics** – Available at `/metrics` endpoint (enabled by default). Follows GenAI semantic conventions.
- **OTLP export** – Configure `OTLP_ENDPOINT` to push traces/metrics to a collector.
- **OpenTelemetry** – Integrated for distributed tracing.

## MCP (Model Context Protocol) Integration

- MCP servers can be registered with the gateway.
- Tools exposed by MCP servers are automatically discovered and made available to LLMs via function calling.
- This enables dynamic tool invocation without extra configuration.

## Agent-to-Agent (A2A) & Agent Defineition Language (ADL)

- A2A allows multiple agents to coordinate and communicate.
- ADL lets you define agents as code; CLI can generate Go or Rust servers from ADL definitions.
- Ideal for building multi‑agent systems on top of the gateway.

## Middleware Control

- You can enable/disable specific middleware (rate limiting, logging, CORS, etc.) via environment flags or configuration files.
- Bypass routes can be defined for health checks or internal endpoints.

## Useful Commands (CLI)

If you have the CLI tool installed:

```bash
ig-cli health                 # Check gateway status
ig-cli chat --model openai/gpt-5  # Interactive chat
ig-cli config apply -f config.yaml  # Apply YAML configuration
```

---

For full details, see the [official documentation](https://docs.inference-gateway.com).


Feel free to copy this as a third file, combine it with the others, or ask for a different focus. I’m happy to adjust further.
