# RFC-005 — API Gateway

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [RFC-003](./RFC-003-estrategia-de-autenticacao.md), [ADR-007](../adr/ADR-007-lambda-authorizer-na-borda.md) |

## Resumo

Usar o **Amazon API Gateway HTTP API (v2)** como porta de entrada única, com integração `HTTP_PROXY` para o NLB do cluster, integração Lambda para `POST /auth/cpf` e um Lambda authorizer nas rotas protegidas. Provisionado em Terraform no repositório `oficina-lambda-auth`.

## Contexto / Problema

O PDF exige um API Gateway para controle e roteamento e a proteção de rotas sensíveis por autenticação via CPF. Hoje o NLB do Kubernetes é exposto diretamente. É necessário um ponto de entrada que (a) encaminhe para o cluster, (b) invoque a Lambda de autenticação e (c) rejeite chamadas sem token antes de chegarem à aplicação — tudo sem adicionar carga aos nós `t3.small` e destruível junto com o resto.

## Opções consideradas

| Critério | API Gateway HTTP API | API Gateway REST API | Kong (no cluster) | Traefik (Ingress no cluster) |
|---|---|---|---|---|
| Integração nativa com Lambda | Sim | Sim | Plugin/serverless-functions | Não nativa |
| Authorizer | Lambda (REQUEST) ou JWT/JWKS | Lambda, Cognito | Plugins JWT (exige RS256 ou segredo no Kong) | Middleware ForwardAuth para serviço externo |
| Proxy para backend HTTP público (NLB) | `HTTP_PROXY` com `overwrite:path` | HTTP proxy | Upstream | Service |
| Custo | US$ 1,00 / milhão de req; free tier 1M/mês no 1º ano | US$ 3,50 / milhão | Pods extras (CPU/RAM) + DB do Kong | Pod extra, leve |
| Recursos consumidos no cluster | **Zero** | Zero | Alto para `t3.small` | Baixo |
| Provisionamento | Terraform, destruível com `destroy` | Terraform | Helm + CRDs | Helm + CRDs |
| Esforço | Baixo | Médio (stages, deployments, mais recursos) | Alto | Médio |
| Observabilidade | Access logs no CloudWatch | idem + métricas detalhadas | Plugins | Métricas Prometheus |
| Risco | Baixo | Baixo | Médio | Médio (não cobre Lambda) |

## Proposta

1. **HTTP API** com stage `$default` e auto-deploy.
2. **Rotas:**
   - `POST /auth/cpf` → integração `AWS_PROXY` com a Lambda `auth-cpf`.
   - Públicas para o cluster: `GET /health`, `GET /swagger/{proxy+}`, `POST /api/v1/auth/login`, `POST /api/v1/auth/register`, `GET /api/v1/service-orders/{id}/status`, `POST /api/v1/service-orders/{id}/budget-approval`.
   - `ANY /api/v1/{proxy+}` → cluster, protegida pelo **Lambda authorizer**.
3. **Integração com o cluster:** `HTTP_PROXY` para `http://<dns-do-nlb>` com `overwrite:path = $request.path`, preservando método, path, query e headers. O DNS do NLB é descoberto pelo Terraform via `data "aws_lbs"` filtrando pela tag `kubernetes.io/service-name = oficina-api/oficina-api`; há uma variável de override (`backend_url`) para ambientes sem descoberta.
4. **Authorizer:** tipo REQUEST, payload 2.0, `enable_simple_responses`, `identity_sources = ["$request.header.Authorization"]`, TTL de 300 s.
5. **Logs de acesso** do gateway no CloudWatch com `requestId`, rota, status e latência de integração.
6. **CORS** liberado para `GET/POST/PUT/PATCH/DELETE` com header `Authorization`, para uso do Swagger e de clientes web.

## Consequências

- **Positivas:** ponto único de entrada e de política de acesso; 401 antes de tocar o cluster; nenhum pod adicional; a URL do gateway é estável enquanto ele existir, independente de recriações do NLB; custo praticamente zero no volume da demonstração.
- **Negativas:** NLB continua público (o gateway fala com ele pela internet) — mitigável com VPC Link para um NLB interno em uma fase futura; a lista de rotas públicas precisa ser mantida em dois lugares (gateway e aplicação); HTTP API não oferece caching nem throttling por API key como o REST API.
- **Operação:** o gateway depende do `Service` do Kubernetes existir; a ordem de criação é k8s → db → app → lambda (ver ADR-004).

## Perguntas abertas

- Migrar para VPC Link + NLB interno para fechar o acesso direto ao cluster? Recomendado para produção.
- Throttling por rota no stage `$default`? O padrão da conta (10k rps) é mais do que suficiente aqui.
