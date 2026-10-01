
# Inference Gateway – Provider Configuration Reference

This file summarises the environment variables and model name prefixes required for each supported LLM provider.

---

## Provider Environment Variables

| Provider | Required Env Var | Optional Env Var | Model Name Preifix |
|----------|------------------|------------------|-------------------|
| **OpenAI** | `OPENAI_API_KEY` | `OPENAI_BASE_URL` | `openai/` (e.g., `openai/gpt-5`) |
| **Anthropic** | `ANTHROPIC_API_KEY` | `ANTHROPIC_BASE_URL` | `anthropic/` (e.g., `anthropic/claude-opus-4-8`) |
| **Groq** | `GROQ_API_KEY` | `GROQ_BASE_URL` | `groq/` (e.g., `groq/llama-3.3-70b-versatile`) |
| **Cohere** | `COHERE_API_KEY` | `COHERE_BASE_URL` | `cohere/` (e.g., `cohere/command-r-plus`) |
| **Ollama** | – (no key) | `OLLAMA_BASE_URL` (default `http://localhost:11434`) | `ollama/` (e.g., `ollama/llama3.2`) |
| **llama.cpp** | – (no key) | `LLAMACPP_BASE_URL` (default `http://localhost:8080`) | `llamacpp/` (e.g., `llamacpp/llama-3.2-3b-instruct`) |
| **Cloudflare** | `CLOUDFLARE_API_KEY` | `CLOUDFLARE_BASE_URL` | `cloudflare/` |
| **DeepSeek** | `DEEPSEEK_API_KEY` | `DEEPSEEK_BASE_URL` | `deepseek/` (e.g., `deepseek/deepseek-v4-flash`) |
| **Google** | `GOOGLE_API_KEY` | `GOOGLE_BASE_URL` | `google/` (e.g., `google/gemini-pro`) |
| **Mistral** | `MISTRAL_API_KEY` | `MISTRAL_BASE_URL` | `mistral/` |
| **MiniMax** | `MINIMAX_API_KEY` | `MINIMAX_BASE_URL` | `minimax/` |
| **Moonshot** | `MOONSHOT_API_KEY` | `MOONSHOT_BASE_URL` | `moonshot/` |
| **Nvidia** | `NVIDIA_API_KEY` | `NVIDIA_BASE_URL` | `nvidia/` |
| **Z.ai** | `ZAI_API_KEY` | `ZAI_BASE_URL` | `zai/` |

---

## Model Name Format

When sending a request to `/v1/chat/completions`, include the provider preifix followed by the actual model ID recognised by that provider.

**Examples:**

- `openai/gpt-5`
- `anthropic/claude-opus-4-8`
- `groq/llama-3.3-70b-versatile`
- `deepseek/deepseek-v4-flash`
- `ollama/llama3.2`

---

## Multiple Providers

You can set multiple API keys simultaneously – the gateway will route requests based on the model preifix.

```bash
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
export GROQ_API_KEY=gsk_...
```

---

## Base URL Overrides

If a provider offers a custom endpoint (e.g., Azure OpenAI, local proxy), set the corresponding `_BASE_URL` variable:

```bash
export OPENAI_BASE_URL=https://your-azure-openai.openai.azure.com/
```

---

## No‑Key Providers

- **Ollama** and **llama.cpp** do not require API keys. They only need a reachable base URL.

---

For the full list of supported providers and model IDs, call:

```bash
curl http://localhost:8080/v1/models
```

This returns all currently configured models with their full names (prefix + ID).
