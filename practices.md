
# Inference Gateway – Production Deployment Best Practices

This guide covers key considerations for deploying Inference Gateway in a production environment, including Kubernetes and containerised setups.

---

## 1. Resource Limits (CPU / Memory)

Set explicit resource requests and limits to avoid noisy neighbours and ensure predictable performance.

**Example (Kubernetes):**

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

- Start with the above and adjust based on your workload (concurrency, model size, etc.).
- Memory usage grows with request size and streaming buffers – monitor and tune.

---

## 2. Health Checks & Readiness/Liveneess Probes

Inference Gateway provides a built‑in health endpoint at `/health`. Use it for Kubernetes probes.

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

- **Readiness** – prevents traffic before the gateway is fully initialised.
- **Liveness** – restarts the pod if it becomes unresponsive.

---

## 3. Horizontal Scaling (HPA)

Deploy multiple replicas and use a HorizontalPodAutoscaler based on CPU or custom metrics (e.g., request rate).

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: inference-gateway
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: inference-gateway
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

- For production, run at least 2 replicas for high availability.
- If using the [Inference Gateway Operator](https://docs.inference-gateway.com/operator/), scaling can be automated declaratively.

---

## 4. Logging Strategy

Enable structured logging (JSON) to integrate with your log aggreagator (ELK, Loki, etc.).

Set environment variable:
```
LOG_FORMAT=json
LOG_LEVEL=info   # use 'warn' in production to reduce noise
```

- Forward logs to a central system for analysis.
- Avoid logging sensitive data (API keys) – the gateway already redacts them by default.

---

## 5. Monitoring & Observability

- **Prometheus metrics** – exposed at `/metrics`. Scrape this endpoint with Prometheus.
- **Grafana dashboards** – import or create dashboards to visualise request rates, latency, errors, and token usage.
- **OTLP export** – set `OTLP_ENDPOINT` to push traces and metrics to a collector (e.g., Tempo, Datadog).

Example Prometheus scrape config:
```yaml
scrape_configs:
  - job_name: 'inference-gateway'
    static_configs:
      - targets: ['inference-gateway:8080']
```

---

## 6. Security – TLS & Authentication

### TLS (HTTPS)

Enable TLS to encrypt traffic:
```
TLS_ENABLE=true
TLS_CERT_FILE=/path/to/tls.crt
TLS_KEY_FILE=/path/to/tls.key
```

Use cert-manager with Kubernetes to automate certificate renewal.

### Authentication (OIDC)

Enable OIDC for production:
```
AUTH_ENABLE=true
OIDC_ISSUER=https://your-idp.example.com
OIDC_AUDIENCE=your-api-audience
```

All requests must then include a valid JWT in the `Authorization` header. Ensure your IdP is highly available.

---

## 7. Configuration Management

Store environment varieables and secrets securely.

- Use **Kubernetes Secrets** for API keys:
  ```yaml
  apiVersion: v1
  kind: Secret
  metadata:
    name: llm-secrets
  type: Opaque
  data:
    openai-api-key: <base64-encoded>
  ```

- Use **ConfigMaps** for non‑sensitive settings (e.g., log level, timeouts).

- Consider using the [Operator](https://github.com/inference-gateway/operator) to manage the entire configuration as a custom resource.

---

## 8. Upstream Timeouts & Retries

Set reasonable timeouts for upstream LLM providers to avoid hanging requests.

```
REQUEST_TIMEOUT=60   # seconds
```

- For slow models, increase this value.
- Implement retries with exponential backoff at the client level if needed (the gateway does not retry by default).

---

## 9. Persistent Storage (If Needed)

Inference Gateway is stateless – it does not require persistent storage. This simplifies scaling and disaster recovery.

If you use caching (e.g., Redis for response caching), manage that as an external dependency.

---

## 10. Backup & Disaster Recovery

Since the gateway is stateless, recovery is straightforward:
- Re‑deploy the same configuration.
- Ensure all secrets and ConfigMaps are backed up (use GitOps with ArgoCD or Flux).
- Regularly back up your IdP configuration (if you manage custom users).

---

## Summary Checklist

- [ ] Resource requests/limits defined
- [ ] Health probes configured
- [ ] HPA set up (≥2 replicas)
- [ ] Structured logging enabled
- [ ] Prometheus scraping configured
- [ ] TLS enabled (if exposed publicly)
- [ ] OIDC authentication enabled (if needed)
- [ ] Secrets stored securely
- [ ] Timeouts tuned for your LLM providers
- [ ] Configuration under version control (GitOps)

---

For more details, refer to the [official documentation](https://docs.inference-gateway.com).


Feel free to copy, adjust, and combine with your existing files. Let me know if you need a different angle next time.
