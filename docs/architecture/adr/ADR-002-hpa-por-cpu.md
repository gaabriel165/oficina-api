# ADR-002 — Autoescala horizontal por CPU (HPA)

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## Contexto

A aplicação roda em um cluster EKS com node group de 2 a 4 nós `t3.small` (2 vCPU, 2 GiB). A direção da oficina exige escalabilidade e alta disponibilidade. A aplicação é *stateless* (estado no RDS), então escalar horizontalmente é seguro. Precisávamos escolher o mecanismo e a métrica de escala.

## Decisão

Usar o **HorizontalPodAutoscaler** (`autoscaling/v2`) sobre o `Deployment` da aplicação:

- `minReplicas: 2` (alta disponibilidade mínima — um pod por nó em condições normais) e `maxReplicas: 6`.
- Métrica: **utilização média de CPU a 60%** dos `requests` (`100m` por pod; `limits 500m`).
- `metrics-server` instalado pelo CD do repositório `oficina-infra-k8s` como pré-requisito.
- Probes `startup`, `readiness` e `liveness` em `/health` para que réplicas novas só recebam tráfego quando prontas.
- Escala de **nós** delegada ao managed node group (2..4); não há Cluster Autoscaler/Karpenter nesta fase.

Em homologação o `Deployment` roda com 1 réplica e sem HPA para economizar recursos.

## Consequências

**Positivas**
- Demonstrável em minutos com um teste de carga (`k6` em `load/`): as réplicas sobem de 2 para 4–6 e voltam.
- Métrica simples, disponível sem instrumentação adicional; correlaciona bem com a carga de requisições HTTP em uma API Go.
- Baixo `requests.cpu` permite empacotar várias réplicas por nó pequeno.

**Negativas**
- CPU não captura gargalos de I/O (banco lento) — a aplicação pode escalar sem que o RDS acompanhe. Mitigação: métricas de latência e Performance Insights no monitoramento.
- Teto de 6 réplicas limitado pela capacidade dos nós; com 4 nós `t3.small` sobra pouco para o agente do New Relic. Ajustável por variável.
- Cooldown padrão do HPA (5 min para reduzir) mantém réplicas extras por alguns minutos após o pico.

## Alternativas rejeitadas

- **Escalar por requisições por segundo (métrica customizada via Prometheus Adapter ou KEDA)** — mais fiel à carga, mas exige componentes extras no cluster e configuração que não cabe no prazo.
- **Escalar por memória** — a aplicação Go tem uso de memória estável; não reflete carga.
- **Réplicas fixas** — não atende ao requisito de escalabilidade nem à demonstração.
