# Documentação de Arquitetura — Tech Challenge Fase 3

| | |
|---|---|
| **Projeto** | Oficina Mecânica — Sistema Integrado de Atendimento e Execução de Serviços |
| **Fase** | 3 — Segurança, serverless, CI/CD segregado e observabilidade |
| **Status** | Vigente |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

Esta pasta reúne a documentação arquitetural completa exigida na Fase 3. Ela foi organizada em níveis de zoom, do mais amplo ao mais específico, inspirada no modelo C4: primeiro **onde** o sistema vive e com quem conversa, depois **como** as peças se comunicam, em seguida **quais dados** persistimos e, por fim, **por que** cada decisão foi tomada.

## Como ler

1. **Contexto e componentes** → [`components.md`](./components.md)
   Visão de nuvem com API Gateway, Lambdas, EKS, RDS, New Relic e os quatro repositórios. Inclui também as camadas internas da aplicação (Clean Architecture).
2. **Fluxos (diagramas de sequência)**
   - [`sequence-authentication.md`](./sequence-authentication.md) — autenticação de cliente por CPF (Lambda) e de operador (login), validação do token no gateway e na aplicação.
   - [`sequence-service-order.md`](./sequence-service-order.md) — abertura de ordem de serviço em chamada única e transições de status com notificação e telemetria.
3. **Dados** → [`database.md`](./database.md)
   Justificativa formal do PostgreSQL, diagrama ER, explicação dos relacionamentos e os ajustes do modelo relacional feitos nesta fase (status do cliente, índices, normalização).
4. **Decisões**
   - **RFCs** (decisões técnicas discutíveis, com opções comparadas) em [`rfc/`](./rfc/).
   - **ADRs** (decisões arquiteturais permanentes) em [`adr/`](./adr/).
5. **Diagramas por repositório** → [`repo-diagrams.md`](./repo-diagrams.md)
   Um diagrama pequeno para o README de cada um dos quatro repositórios.

## RFCs — Request for Comments

| # | Título | Status |
|---|---|---|
| [RFC-001](./rfc/RFC-001-escolha-da-nuvem.md) | Escolha da nuvem | Aceita |
| [RFC-002](./rfc/RFC-002-escolha-do-banco-de-dados.md) | Escolha do banco de dados | Aceita |
| [RFC-003](./rfc/RFC-003-estrategia-de-autenticacao.md) | Estratégia de autenticação | Aceita |
| [RFC-004](./rfc/RFC-004-ferramenta-de-observabilidade.md) | Ferramenta de observabilidade | Aceita |
| [RFC-005](./rfc/RFC-005-api-gateway.md) | API Gateway | Aceita |

## ADRs — Architecture Decision Records

| # | Título | Status |
|---|---|---|
| [ADR-001](./adr/ADR-001-comunicacao-sincrona-rest.md) | Comunicação síncrona via REST | Aceita |
| [ADR-002](./adr/ADR-002-hpa-por-cpu.md) | Autoescala horizontal por CPU (HPA) | Aceita |
| [ADR-003](./adr/ADR-003-jwt-hs256-segredo-compartilhado.md) | JWT HS256 com segredo compartilhado via SSM | Aceita |
| [ADR-004](./adr/ADR-004-quatro-repositorios-e-remote-state.md) | Quatro repositórios e remote state no S3 | Aceita |
| [ADR-005](./adr/ADR-005-ambientes-por-namespace.md) | Homologação e produção como namespaces | Aceita |
| [ADR-006](./adr/ADR-006-logs-estruturados-com-correlacao.md) | Logs estruturados em JSON com correlação | Aceita |
| [ADR-007](./adr/ADR-007-lambda-authorizer-na-borda.md) | Lambda authorizer na borda | Aceita |
| [ADR-008](./adr/ADR-008-migrations-no-startup-da-aplicacao.md) | Migrations executadas no startup da aplicação | Aceita |

## Os quatro repositórios

| Repositório | Responsabilidade | Provisiona / entrega |
|---|---|---|
| [`oficina-api`](https://github.com/gaabriel165/oficina-api) | Aplicação principal (Go, Clean Architecture) | Imagem Docker no ECR, deploy em EKS via kustomize (namespaces `oficina-api-homolog` e `oficina-api`) |
| [`oficina-infra-k8s`](https://github.com/gaabriel165/oficina-infra-k8s) | Infraestrutura Kubernetes (Terraform) | VPC, EKS 1.34, ECR, OIDC + roles do GitHub Actions, segredos no SSM, metrics-server e New Relic via Helm |
| [`oficina-infra-db`](https://github.com/gaabriel165/oficina-infra-db) | Banco de dados gerenciado (Terraform) | RDS PostgreSQL 16 privado, security groups, `DATABASE_URL` no SSM |
| [`oficina-lambda-auth`](https://github.com/gaabriel165/oficina-lambda-auth) | Function serverless (Go) + API Gateway (Terraform) | Lambdas `auth-cpf` e `authorizer`, API Gateway HTTP API roteando para o cluster |

## Convenções

- Diagramas em **Mermaid**, renderizados nativamente pelo GitHub. Versões em PNG são geradas para o documento de entrega.
- Datas no formato ISO (`AAAA-MM-DD`).
- Nomes de recursos AWS prefixados por `oficina-api`; parâmetros no SSM sob `/oficina-api/`.
