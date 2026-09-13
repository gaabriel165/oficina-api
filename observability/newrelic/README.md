# Observabilidade — New Relic

Assets de monitoramento versionados como código. Aplicados uma vez por conta com o script abaixo; os dados chegam por três caminhos:

| Sinal | Origem | Como chega |
|---|---|---|
| Logs JSON (`http.request`, `service_order.*`, `integration.error`) | stdout dos pods | Fluent Bit do `nri-bundle` (instalado pelo CD do `oficina-infra-k8s`) → Logs, com os atributos JSON extraídos |
| Traces (spans HTTP e `gorm.*`) | aplicação (`otelgin` + plugin GORM) | OTLP/HTTP direto para `otlp.nr-data.net:4318` (`OTEL_EXPORTER_OTLP_HEADERS=api-key=<license>`) |
| CPU, memória, HPA, eventos | `nri-bundle` (infra agent, kube-state-metrics) | Kubernetes integration |
| Uptime | Synthetics ping monitor `oficina-api-health` | Criado pelo script |

## Aplicar

```bash
export NEW_RELIC_API_KEY=NRAK-...          # User API key
export NEW_RELIC_ACCOUNT_ID=1234567
export ALERT_EMAIL=voce@exemplo.com
./observability/newrelic/setup.sh

# depois que o API Gateway existir, só o monitor de uptime:
ONLY_SYNTHETICS=true HEALTH_URL=https://<api-id>.execute-api.us-east-1.amazonaws.com/health \
  ./observability/newrelic/setup.sh
```

`HEALTH_URL` é opcional na primeira execução: sem ela o monitor Synthetics é pulado e pode ser criado depois com `ONLY_SYNTHETICS=true`.

Cria:

- **Dashboard "Oficina API — Operação"** (`dashboard.json`): volume diário de OS, tempo médio por status, falhas de OS, erros de integração, latência p50/p95/p99 e por rota, throughput, status HTTP, tempo em banco, CPU/memória por pod, réplicas do HPA, CPU dos nós e uptime.
- **Política de alertas "Oficina API"** com condições NRQL (agregação `CADENCE`, janela de 60 s e atraso de 120 s — necessária porque eventos como `service_order.failed` são esparsos e, com `EVENT_FLOW`, a janela só fecharia quando chegasse o próximo evento):
  - falha no processamento de OS (`service_order.failed` > 0 em 5 min);
  - erros de integração (`integration.error` > 2 em 5 min);
  - erros 5xx na API;
  - latência p95 acima de 1 s.
- **Notificação por e-mail** (destination + channel + workflow ligado à política).
- **Monitor Synthetics** de ping no `/health` a cada 5 minutos.

## Consultas NRQL de referência

```sql
-- volume diário de OS
SELECT count(*) FROM Log WHERE event = 'service_order.created' FACET dateOf(timestamp) SINCE 30 days ago
-- tempo médio em cada status
SELECT average(numeric(duration_minutes)) FROM Log WHERE event = 'service_order.status_changed' FACET from_status
-- erros nas integrações
SELECT count(*) FROM Log WHERE event = 'integration.error' FACET integration, operation TIMESERIES
-- correlação de uma requisição
SELECT * FROM Log WHERE request_id = '<X-Request-ID>' ORDER BY timestamp
```
