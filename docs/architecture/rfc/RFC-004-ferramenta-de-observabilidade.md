# RFC-004 — Ferramenta de observabilidade

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [ADR-006](../adr/ADR-006-logs-estruturados-com-correlacao.md), [sequence-service-order.md](../sequence-service-order.md) |

## Resumo

Adotar o **New Relic** (plano gratuito permanente) para traces, logs, métricas de Kubernetes, dashboards, alertas e verificação de uptime, alimentado por **OpenTelemetry** na aplicação e pelo bundle `nri-bundle` no cluster.

## Contexto / Problema

O PDF pede integração com Datadog, New Relic ou equivalente, monitorando latência das APIs, CPU/memória do Kubernetes, healthchecks/uptime, alertas de falha no processamento de OS e logs JSON com correlação, além de dashboards de volume diário de OS, tempo médio por status e erros de integração. O ambiente é efêmero e o avaliador verá os dashboards **depois** que a infra for destruída, então a ferramenta precisa reter dados sem custo e sem prazo de trial.

## Opções consideradas

| Critério | New Relic (Free) | Datadog (Trial) | Grafana Cloud (Free) | Prometheus + Grafana no cluster |
|---|---|---|---|---|
| Custo | **US$ 0 permanente**: 100 GB/mês de ingestão, 1 usuário full, retenção 8 dias | Gratuito só por **14 dias**; depois host-based (≈ US$ 15/host/mês + logs) | US$ 0: 10k séries, 50 GB logs, 50 GB traces, 14 dias | US$ 0 de licença, mas consome CPU/RAM dos nós `t3.small` |
| Dashboards visíveis após destruir a infra | Sim (8 dias) | Só se a trial não expirou | Sim (14 dias) | **Não** — morrem com o cluster |
| Traces OTLP nativo | Sim (`otlp.nr-data.net:4318`) | Sim (via agent) | Sim (Tempo) | Precisa de Tempo/Jaeger extra |
| Logs de pods com parse de JSON automático | Sim (`newrelic-logging` / Fluent Bit) | Sim (agent) | Loki + Alloy | Loki extra |
| Métricas Kubernetes (CPU, memória) prontas | `nri-bundle` + kube-state-metrics, dashboards prontos | Agent + dashboards prontos | Alloy + dashboards | kube-prometheus-stack (pesado) |
| Alertas por e-mail | NRQL alert conditions | Monitors | Alerting | Alertmanager |
| Uptime externo | Synthetics ping (100 checks/mês grátis) | Synthetics pago | Synthetic Monitoring free tier | Blackbox exporter |
| Esforço de instalação | Helm chart + license key; OTLP direto do app | Helm chart + API key | Helm chart + 3 tokens | Vários charts; tuning de recursos |
| Risco para a demonstração | Baixo | **Alto** — trial pode expirar antes da avaliação | Médio (menos integrado) | Alto — nós pequenos, sem persistência após destroy |

## Proposta

1. **Aplicação:** `log/slog` com handler JSON para stdout; middleware de `X-Request-ID`; OpenTelemetry (`otelgin`) exportando spans via OTLP/HTTP para o New Relic com `service.name = oficina-api`; `trace_id`, `span_id` e `request_id` em todo log; eventos de negócio `service_order.created`, `service_order.status_changed`, `service_order.failed`, `integration.error`.
2. **Cluster:** CD do `oficina-infra-k8s` instala `nri-bundle` (infra agent, kube-state-metrics, Fluent Bit para logs, kube-events) via Helm com a license key de um GitHub Secret. Também instala `metrics-server` (pré-requisito do HPA).
3. **Dashboards NRQL** (criados via NerdGraph): volume diário de OS; tempo médio em cada status (`FACET from_status`); erros de integração por `integration`; latência p50/p95 por rota (spans) e por `latency_ms` dos logs; CPU e memória por pod (`K8sContainerSample`); taxa de 5xx.
4. **Alerta:** condição NRQL `count(*) WHERE event = 'service_order.failed' > 0` em janela de 5 min, notificação por e-mail.
5. **Uptime:** monitor Synthetics tipo *ping* em `GET /health` via API Gateway, a cada 5 min.
6. **Lambdas e RDS** continuam com logs no CloudWatch (padrão AWS); não são enviados ao New Relic nesta fase para caber na ingestão gratuita.

## Consequências

- **Positivas:** todos os itens do PDF cobertos por uma única ferramenta sem custo; dashboards sobrevivem ao `terraform destroy` por 8 dias; instrumentação por OpenTelemetry é portável — trocar o backend é alterar endpoint e header.
- **Negativas:** 1 único usuário full no plano gratuito (o avaliador vê dashboards por link público/compartilhado ou pelo vídeo); retenção de 8 dias exige gravar o vídeo com dados recentes; o agente de infraestrutura consome ~150 MB de RAM por nó.
- **Riscos mitigados:** se a ingestão passar de 100 GB/mês (improvável), o New Relic apenas para de aceitar dados — não gera cobrança sem upgrade explícito.

## Perguntas abertas

- Enviar logs das Lambdas ao New Relic via extensão Lambda? Útil para correlacionar a autenticação com a requisição; adiado.
- Substituir eventos em log por métricas OTel (counters/histograms)? Reduz custo de ingestão em escala; para o volume atual os logs são suficientes e mais fáceis de investigar.
