# Diagrama de Sequência — Ordem de Serviço

| | |
|---|---|
| **Status** | Vigente |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## 1. Abertura de OS em chamada única

Desde a Fase 2 a OS é aberta com cliente, veículo, serviços e peças no mesmo `POST`. Na Fase 3 a requisição passa pelo API Gateway e pelo authorizer, e a aplicação passa a emitir o evento `service_order.created` nos logs estruturados, que alimenta o dashboard de volume diário no New Relic.

```mermaid
sequenceDiagram
    autonumber
    actor Operador
    participant GW as API Gateway HTTP API
    participant AZ as Lambda authorizer
    participant H as ServiceOrderHandler
    participant UC as CreateServiceOrderUseCase
    participant CR as CustomerRepository
    participant VR as VehicleRepository
    participant SR as ServiceRepository
    participant PR as PartRepository
    participant OR as ServiceOrderRepository
    participant DB as RDS PostgreSQL
    participant NR as New Relic

    Operador->>GW: POST /api/v1/service-orders { customer_id, vehicle_id, notes, services[], parts[] }
    GW->>AZ: valida Bearer JWT
    AZ-->>GW: isAuthorized = true
    GW->>H: proxy via NLB (X-Request-ID gerado ou propagado)
    H->>H: bind e validação do DTO
    H->>UC: Execute(input)

    UC->>CR: FindByID(customer_id)
    CR->>DB: SELECT customers
    DB-->>CR: cliente
    CR-->>UC: Customer

    UC->>VR: FindByID(vehicle_id)
    VR->>DB: SELECT vehicles
    DB-->>VR: veículo
    VR-->>UC: Vehicle
    UC->>UC: veículo pertence ao cliente? senão ErrServiceOrderVehicleNotOwnedByCustomer

    UC->>UC: entity.NewServiceOrder(status = received)

    loop para cada service_id
        UC->>SR: FindByID(service_id)
        SR->>DB: SELECT services
        DB-->>SR: serviço
        UC->>UC: order.AddItem(service)  snapshot de nome e preço
    end

    loop para cada part
        UC->>PR: FindByID(part_id)
        PR->>DB: SELECT parts
        DB-->>PR: peça
        UC->>UC: HasSufficientStock? senão ErrPartInsufficientStock
        UC->>UC: order.AddPart(part, quantity)  snapshot de nome e preço
    end

    UC->>OR: Create(order)
    OR->>DB: transação - INSERT service_orders, service_order_items, service_order_parts
    DB-->>OR: ok
    OR-->>UC: nil
    UC-->>H: ServiceOrder

    H->>NR: log JSON event=service_order.created order_id status=received items parts request_id trace_id
    H-->>Operador: 201 Created { id, status, total_amount, items, parts }

    Note over H,NR: Em caso de erro em qualquer etapa o handler mapeia para 404/422 e emite event=service_order.failed com o motivo
```

### Tratamento de erros nesta rota

| Erro de domínio | HTTP | Evento emitido |
|---|---|---|
| `ErrCustomerNotFound`, `ErrVehicleNotFound`, `ErrServiceNotFound`, `ErrPartNotFound` | 404 | `service_order.failed` |
| `ErrServiceOrderVehicleNotOwnedByCustomer`, `ErrPartInsufficientStock`, `ErrServiceOrderItemInvalidQty` | 422 | `service_order.failed` |
| erro de infraestrutura (banco) | 500 sem vazar detalhe | `service_order.failed` |

## 2. Transição de status com notificação e telemetria

Todas as transições (`start-diagnosis`, `send-budget`, `approve-budget`, `reject-budget`, `finish-execution`, `deliver` e o webhook `budget-approval`) seguem o mesmo padrão: a entidade valida a máquina de estados, o repositório persiste, o notificador envia e-mail em modo best-effort e o evento `service_order.status_changed` registra quanto tempo a OS ficou no status anterior.

```mermaid
sequenceDiagram
    autonumber
    actor Operador
    participant GW as API Gateway HTTP API
    participant H as ServiceOrderHandler
    participant UC as ApproveBudgetUseCase
    participant OR as ServiceOrderRepository
    participant PR as PartRepository
    participant DB as RDS PostgreSQL
    participant N as OrderStatusNotifier
    participant CR as CustomerRepository
    participant RS as Resend
    participant NR as New Relic

    Operador->>GW: PATCH /api/v1/service-orders/{id}/approve-budget  Bearer JWT
    GW->>H: authorizer ok, proxy via NLB
    H->>UC: Execute(id)
    UC->>OR: FindByID(id)
    OR->>DB: SELECT service_orders + items + parts
    DB-->>OR: OS em waiting_approval
    OR-->>UC: ServiceOrder

    UC->>UC: order.ApproveBudget()  waiting_approval -> in_execution, started_at = now
    Note right of UC: transição inválida devolve ErrServiceOrderInvalidTransition (422)

    loop para cada peça da OS
        UC->>PR: FindByID(part_id)
        PR->>DB: SELECT parts
        UC->>UC: part.DebitStock(quantity)
        UC->>PR: Update(part)
        PR->>DB: UPDATE parts SET stock_quantity
    end

    UC->>OR: Update(order)
    OR->>DB: UPDATE service_orders SET status, started_at, updated_at
    DB-->>OR: ok

    UC->>N: Notify(order)
    N->>CR: FindByID(customer_id)
    CR->>DB: SELECT customers
    DB-->>CR: cliente
    N->>RS: POST /emails  (assunto com status, HTML e texto)
    alt Resend falhou
        RS-->>N: erro
        N->>NR: log JSON event=integration.error integration=resend order_id error
        Note over N: best-effort: a transição já foi persistida, a falha de e-mail não é propagada
    else enviado
        RS-->>N: 200
    end

    UC-->>H: ServiceOrder (status = in_execution)
    H->>NR: log JSON event=service_order.status_changed from_status=waiting_approval to_status=in_execution duration_minutes
    H-->>Operador: 200 OK { id, status, total_amount, started_at }
```

### Como o New Relic usa esses eventos

| Requisito do PDF | Evento / fonte | Consulta NRQL (exemplo) |
|---|---|---|
| Volume diário de ordens de serviço | `service_order.created` | `SELECT count(*) FROM Log WHERE event = 'service_order.created' TIMESERIES 1 day` |
| Tempo médio por status (diagnóstico, execução, finalização) | `service_order.status_changed` com `from_status` e `duration_minutes` | `SELECT average(duration_minutes) FROM Log WHERE event = 'service_order.status_changed' FACET from_status` |
| Erros e falhas nas integrações | `integration.error` | `SELECT count(*) FROM Log WHERE event = 'integration.error' FACET integration` |
| Alertas para falhas no processamento de OS | `service_order.failed` | condição NRQL `SELECT count(*) FROM Log WHERE event = 'service_order.failed'` acima de 0 em 5 min |
| Latência das APIs | spans `otelgin` e campo `latency_ms` do log de requisição | `SELECT percentile(duration.ms, 50, 95) FROM Span WHERE service.name = 'oficina-api' FACET http.route` |

Todos os logs carregam `request_id`, `trace_id` e `span_id`, então a partir de um evento de negócio é possível abrir o trace da requisição e vice-versa (ver [ADR-006](./adr/ADR-006-logs-estruturados-com-correlacao.md)).
