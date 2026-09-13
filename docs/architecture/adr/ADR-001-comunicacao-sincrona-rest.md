# ADR-001 — Comunicação síncrona via REST/HTTP

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## Contexto

O sistema é um monólito modular em Go que expõe APIs para operadores e clientes. Na Fase 3 entram dois novos componentes fora do processo da aplicação — a Lambda de autenticação e o API Gateway — e integrações externas já existentes (Resend para e-mail, BrasilAPI para CNPJ). Era preciso decidir o padrão de comunicação entre esses componentes: síncrono (HTTP/REST) ou assíncrono (filas/eventos como SQS, SNS ou EventBridge).

## Decisão

Toda comunicação entre componentes é **síncrona, via HTTP/REST com JSON**:

- Cliente/Operador → API Gateway → (Lambda `auth-cpf` | NLB → aplicação).
- API Gateway → Lambda `authorizer` (invocação síncrona do authorizer).
- Aplicação → Resend e BrasilAPI (HTTP, com a notificação em modo *best-effort*: falha de e-mail não interrompe a operação nem devolve erro ao chamador).
- Lambda → RDS e aplicação → RDS (conexão direta PostgreSQL).

Não há fila ou barramento de eventos nesta fase. Efeitos colaterais que não precisam bloquear (e-mail) são tratados como chamadas síncronas toleradas a falhas e registradas como `integration.error` nos logs.

## Consequências

**Positivas**
- Fluxos fáceis de depurar: um `request_id` e um `trace_id` percorrem toda a requisição ponta a ponta.
- Nenhuma infraestrutura adicional (SQS, DLQ, consumidores) para provisionar, monitorar e destruir.
- Contratos REST já documentados no Swagger; o API Gateway apenas os expõe.
- Aderência ao requisito de "consumo das APIs protegidas" com respostas imediatas.

**Negativas**
- Latência do e-mail entra no caminho da requisição de transição de status (mitigado: Resend responde em centenas de ms; a chamada não é retentada).
- Sem *retry* nem *replay* de notificações perdidas; uma falha do Resend é apenas registrada.
- Escalabilidade de picos depende do HPA (ADR-002), não de amortecimento por fila.

## Alternativas rejeitadas

- **Eventos de domínio em SQS/SNS para notificações** — melhor resiliência, mas adiciona fila, DLQ e um consumidor (ou segunda Lambda), com custo de operação e de documentação desproporcional ao volume de uma oficina.
- **gRPC entre gateway e aplicação** — API Gateway HTTP API não integra gRPC; ganho de desempenho irrelevante aqui.

## Evolução prevista

Quando surgir integração com múltiplas unidades da oficina ou processamento demorado (ex.: geração de relatórios), a transição de status passa a publicar um evento (SNS/EventBridge) consumido por notificadores; a entidade `ServiceOrder` já expõe a transição (`from`, `to`, duração) de forma a tornar isso uma mudança de adapter, não de domínio.
