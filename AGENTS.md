# AGENTS.md

## Setup commands
- Install deps: `go mod tidy`
- Start dev server: `go run ./cmd/api`
- Run tests: `go test ./...`

## Code style
- gofmt obrigatorio
- golangci-lint
- Erros com contexto (fmt.Errorf)

## Arquitetura
- API /api/v1
- Motor de roteamento em Go
- gRPC interno para servicos auxiliares

## Padrões de logging
- JSON com traceId, routeId, destination

## Estrategia de testes
- Unitarios do algoritmo
- Integracao com Postgres e Redis
- Teste de latencia

## Regras de seguranca
- Rate limiting
- Validacao de input
- Auditoria de mudanca de regras

## Checklist de PR
- gofmt e lint ok
- Tests ok
- Docs atualizadas

## Diretrizes de performance
- p95 /routes < 50ms
- Timeouts curtos em chamadas externas
