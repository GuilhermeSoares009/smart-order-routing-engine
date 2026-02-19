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

- Healthcheck: http://localhost:8082/api/v1/health

## Contratos de API
- OpenAPI: docs/api/openapi.yaml

## Documentacao
- Project Reference Guide: PROJECT_REFERENCE_GUIDE.md
- Especificacoes: docs/specs/spec.md
- Plano tecnico: docs/specs/plan.md
- Tarefas: docs/specs/tasks.md
- ADRs: docs/adr/
- Trade-offs: docs/trade-offs.md
- Threat model: docs/threat-model.md
- Performance budget: docs/performance-budget.md
- Feature flags: docs/feature-flags.md
- Legacy spec (arquivado): docs/legacy-spec/
