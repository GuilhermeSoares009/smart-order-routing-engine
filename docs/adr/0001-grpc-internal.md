# ADR 0001 - gRPC Interno

## Status
Aceito

## Contexto
Precisamos baixa latencia entre servicos internos.

## Decisao
Usar gRPC internamente e REST para clientes externos.

## Consequencias
- Melhor performance interna
- Mais complexidade no contrato
