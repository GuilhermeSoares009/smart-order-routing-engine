# Smart Order Routing Engine

Motor de roteamento de baixa latencia que seleciona o melhor destino com base em sinais de health e latencia.

## Capacidades-chave
- Roteamento dinamico com fallback
- Circuit breaker e padroes de resiliencia
- Logs estruturados e tracing com OpenTelemetry
- Trilha de auditoria para decisoes de roteamento

## Inicio rapido (Docker)
```bash
docker compose up --build
```

- Inclui `docker-compose.yml` e Dockerfile(s).
- Healthcheck: http://localhost:8082/api/v1/health

## Qualidade (pre-commit)
Este repositorio usa pre-commit para CR + auditoria ASVS (OWASP ASVS v5.0.0) antes de cada commit.

```bash
pip install pre-commit
pre-commit install
```

Para rodar manualmente:

```bash
pre-commit run --all-files
```
