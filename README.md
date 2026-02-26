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

## API (MVP)

- `GET /api/v1/health`
- `POST /api/v1/routes`
- `GET /api/v1/audit/routes`

### Variaveis de ambiente

- `PORT` (default: 8080)
- `RATE_LIMIT_PER_MIN` (default: 120)
