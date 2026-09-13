# ADR-006 — Logs estruturados em JSON com correlação de requisições

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [RFC-004](../rfc/RFC-004-ferramenta-de-observabilidade.md) |

## Contexto

Até a Fase 2 a aplicação usava `gin.Default()` (logger de texto) e `log.Printf` nos notificadores. O PDF exige logs estruturados (JSON) com correlação entre requisições, dashboards de negócio (volume de OS, tempo por status, erros de integração) e alertas de falha de processamento. Precisávamos definir formato, campos obrigatórios, onde os eventos de negócio são emitidos e como isso respeita a regra de dependência da Clean Architecture.

## Decisão

1. **Formato:** `log/slog` da biblioteca padrão com `JSONHandler` em stdout, nível `info` por padrão. Campos base em todo registro: `time`, `level`, `msg`, `service = oficina-api`, `environment`, `version`.
2. **Correlação:**
   - Middleware `RequestID`: aceita `X-Request-ID` do chamador (o API Gateway envia o seu) ou gera um UUID; devolve o mesmo header na resposta.
   - OpenTelemetry (`otelgin`) abre um span por requisição; `trace_id` e `span_id` do span corrente são incluídos em todo log emitido dentro da requisição.
   - Logger enriquecido guardado no `context.Context` da requisição; quem loga usa `observability.FromContext(ctx)`.
3. **Log de acesso** por requisição: `method`, `route`, `path`, `status`, `latency_ms`, `client_ip`, `user_agent`, `request_id`, `trace_id`.
4. **Eventos de negócio** como registros com campo `event`:

| `event` | Emitido em | Campos |
|---|---|---|
| `service_order.created` | handler após `CreateServiceOrder` | `order_id`, `customer_id`, `status`, `items`, `parts` |
| `service_order.status_changed` | handler após qualquer transição | `order_id`, `from_status`, `to_status`, `duration_minutes` (tempo no status anterior, exposto pela entidade) |
| `service_order.failed` | handler quando um caso de uso de OS retorna erro | `order_id` (se houver), `operation`, `error`, `http_status` |
| `integration.error` | adapters/notificador em falha de Resend ou BrasilAPI | `integration`, `operation`, `error` |

5. **Onde fica o código:** o adapter de observabilidade vive em `internal/infrastructure/observability` (depende de `slog` e OTel). Handlers e middlewares (camada API) o consomem. Casos de uso e domínio **não** importam logger — a entidade apenas expõe dados (`LastTransition`) e o handler registra. O notificador de status (application) registra a falha da porta `NotificationService` como `integration.error` usando apenas `log/slog` da biblioteca padrão, sem dependência de framework.
6. **Transporte:** stdout → Fluent Bit (`nri-bundle`) → New Relic Logs, que reconhece JSON e indexa cada campo; traces vão direto por OTLP/HTTP.

## Consequências

**Positivas**
- Um `request_id` liga o log do gateway, o log de acesso, os eventos de negócio e o trace; um `trace_id` abre o span no New Relic a partir de qualquer linha de log.
- Dashboards e alertas do PDF são consultas NRQL simples sobre `event`, sem instrumentação de métricas customizadas.
- Sem dependência de framework de log de terceiros; `slog` é padrão desde Go 1.21.
- Regra de dependência preservada: domínio e casos de uso continuam sem logger.

**Negativas**
- Eventos em log custam mais ingestão que métricas agregadas; irrelevante no volume atual.
- A duração por status depende de a entidade registrar a última transição — um campo a mais no agregado.
- Logs de acesso do Gin duplicam parcialmente os spans; mantidos porque a retenção de logs é mais barata e a consulta mais direta.

## Alternativas rejeitadas

- **`zap`/`zerolog`** — mais rápidos, porém desnecessários; `slog` evita dependência.
- **Métricas OTel (counters/histograms) para os eventos de negócio** — melhor para escala, mas o New Relic exigiria dimensões bem planejadas e a investigação de um caso concreto (qual OS falhou?) é mais direta em logs.
- **Emitir eventos dentro dos casos de uso** — acoplaria a camada de aplicação a um logger; preferimos manter a emissão na borda (handlers/adapters).
