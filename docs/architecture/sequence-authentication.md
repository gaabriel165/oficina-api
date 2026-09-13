# Diagrama de Sequência — Autenticação

| | |
|---|---|
| **Status** | Vigente |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

Existem **dois emissores** do mesmo tipo de token (JWT HS256, segredo único em `/oficina-api/jwt_secret`):

| Ator | Rota | Emissor | Credencial | Claims principais |
|---|---|---|---|---|
| Cliente da oficina | `POST /auth/cpf` | Lambda `auth-cpf` | CPF (cliente ativo na base) | `sub` = id do cliente, `cpf`, `name`, `role = customer`, `iss`, `exp` (1 h) |
| Operador | `POST /api/v1/auth/login` | Aplicação (`LoginUseCase`) | e-mail + senha (bcrypt, tabela `users`) | `sub` = id do usuário, `exp` (24 h) |

Ambos os tokens são aceitos pelo **Lambda authorizer** no API Gateway e pelo **middleware JWT** da aplicação. A validação dupla é intencional (defesa em profundidade, ver [ADR-007](./adr/ADR-007-lambda-authorizer-na-borda.md)).

## 1. Autenticação de cliente por CPF

```mermaid
sequenceDiagram
    autonumber
    actor Cliente
    participant GW as API Gateway HTTP API
    participant L as Lambda auth-cpf
    participant SSM as SSM Parameter Store
    participant DB as RDS PostgreSQL

    Note over L,SSM: No cold start a Lambda recebe jwt_secret e database_url como variáveis de ambiente resolvidas pelo Terraform a partir do SSM
    L-->>SSM: leitura na implantação

    Cliente->>GW: POST /auth/cpf { "cpf": "529.982.247-25" }
    GW->>L: invoca (payload HTTP API v2)
    L->>L: normaliza e valida dígitos verificadores do CPF

    alt CPF inválido
        L-->>GW: 400 { "error": "invalid cpf" }
        GW-->>Cliente: 400 Bad Request
    else CPF válido
        L->>DB: SELECT id, name, status FROM customers WHERE document = $1
        alt cliente não encontrado
            DB-->>L: 0 linhas
            L-->>GW: 404 { "error": "customer not found" }
            GW-->>Cliente: 404 Not Found
        else cliente inativo
            DB-->>L: status = inactive
            L-->>GW: 403 { "error": "customer is not active" }
            GW-->>Cliente: 403 Forbidden
        else cliente ativo
            DB-->>L: id, name, status = active
            L->>L: assina JWT HS256 (sub, cpf, name, role=customer, iss, exp = 1h)
            L-->>GW: 200 { "token": "...", "expires_in": 3600, "customer": { id, name } }
            GW-->>Cliente: 200 OK
        end
    end
```

## 2. Consumo de rota protegida com o token

```mermaid
sequenceDiagram
    autonumber
    actor Cliente
    participant GW as API Gateway HTTP API
    participant AZ as Lambda authorizer
    participant NLB as Network Load Balancer
    participant API as oficina-api (EKS)
    participant DB as RDS PostgreSQL

    Cliente->>GW: GET /api/v1/service-orders  Authorization: Bearer <jwt>
    GW->>GW: rota ANY /api/v1/{proxy+} exige authorizer

    alt resultado em cache (mesmo header, até 5 min)
        GW->>GW: reutiliza decisão
    else primeira chamada
        GW->>AZ: REQUEST authorizer (identity source: header Authorization)
        AZ->>AZ: valida prefixo Bearer, assinatura HS256, exp
        alt token ausente, inválido ou expirado
            AZ-->>GW: { isAuthorized: false }
            GW-->>Cliente: 401 Unauthorized
        else token válido
            AZ-->>GW: { isAuthorized: true, context: { sub, role } }
        end
    end

    GW->>NLB: proxy HTTP (path e headers preservados)
    NLB->>API: encaminha ao pod
    API->>API: middleware Request-ID + logger JSON + span OTel
    API->>API: middleware JWT valida HS256 novamente e injeta user_id
    API->>DB: consulta ordens ativas
    DB-->>API: linhas
    API-->>NLB: 200 JSON + X-Request-ID
    NLB-->>GW: 200
    GW-->>Cliente: 200 OK
```

### Rotas que não passam pelo authorizer

| Rota | Motivo |
|---|---|
| `GET /health` | Probe de saúde e uptime (Synthetics) |
| `GET /swagger/{proxy+}` | Documentação pública |
| `POST /api/v1/auth/login`, `POST /api/v1/auth/register` | Emissão de token de operador |
| `POST /auth/cpf` | Emissão de token de cliente (Lambda) |
| `GET /api/v1/service-orders/{id}/status` | Consulta pública de status pelo cliente |
| `POST /api/v1/service-orders/{id}/budget-approval` | Webhook externo autenticado pelo header `X-Webhook-Secret`, validado pela aplicação |

No HTTP API a rota mais específica vence; por isso essas rotas explícitas ficam fora da proteção de `ANY /api/v1/{proxy+}`.

## 3. Login de operador (fluxo existente desde a Fase 1)

```mermaid
sequenceDiagram
    autonumber
    actor Operador
    participant GW as API Gateway HTTP API
    participant API as oficina-api (EKS)
    participant DB as RDS PostgreSQL

    Operador->>GW: POST /api/v1/auth/login { email, password }
    GW->>API: proxy (rota pública)
    API->>DB: SELECT ... FROM users WHERE email = $1
    DB-->>API: password_hash
    API->>API: bcrypt.CompareHashAndPassword
    alt credenciais inválidas
        API-->>Operador: 401 { "error": "invalid credentials" }
    else válidas
        API->>API: assina JWT HS256 (sub = user id, exp = 24h) com o mesmo segredo
        API-->>Operador: 200 { "token": "..." }
    end
```
