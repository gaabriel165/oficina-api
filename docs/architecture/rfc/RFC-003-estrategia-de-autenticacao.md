# RFC-003 — Estratégia de autenticação

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [sequence-authentication.md](../sequence-authentication.md), [ADR-003](../adr/ADR-003-jwt-hs256-segredo-compartilhado.md), [ADR-007](../adr/ADR-007-lambda-authorizer-na-borda.md) |

## Resumo

Autenticar **clientes por CPF** em uma Lambda que consulta o banco e emite **JWT HS256**, com o mesmo segredo já usado pelo login de **operadores** na aplicação. O token é validado na borda por um Lambda authorizer do API Gateway e novamente pelo middleware da aplicação.

## Contexto / Problema

O PDF exige: API Gateway, rotas sensíveis protegidas por autenticação via CPF, e uma function serverless que valida o CPF, consulta existência e status do cliente e devolve um JWT válido para as APIs protegidas. A aplicação já possui autenticação de operadores (e-mail/senha, JWT HS256 de 24 h, middleware Gin). É preciso adicionar a persona **cliente** sem quebrar a existente e sem duplicar lógica de validação.

## Opções consideradas

| Critério | A. Lambda + JWT HS256 (segredo compartilhado) | B. Amazon Cognito (user pool + custom auth) | C. Lambda + JWT RS256 com JWKS público | D. Autenticação só na aplicação (sem Lambda) |
|---|---|---|---|---|
| Atende "function serverless que valida CPF e emite JWT" | Sim | Parcial — emissor é o Cognito; Lambda vira trigger | Sim | **Não** |
| Compatível com o middleware atual (HS256) | **Sim, sem mudança** | Não — tokens RS256 do Cognito exigem JWKS na app | Não — app precisa validar RS256 | Sim |
| Authorizer no API Gateway | Lambda authorizer (REQUEST) | JWT authorizer nativo | JWT authorizer nativo | — |
| Esforço | Baixo | Alto (pool, triggers, fluxo custom para CPF sem senha) | Médio (par de chaves, endpoint JWKS, rotação) | Baixo |
| Custo | Lambda free tier | Cognito free até 10k MAU | Lambda free tier | zero |
| Risco | Segredo em dois lugares (mitigado por SSM único) | Fluxo CPF-sem-senha não é nativo | Rotação de chaves e JWKS a manter | Não cumpre requisito |

## Proposta

1. **Lambda `auth-cpf`** (`POST /auth/cpf`): normaliza e valida os dígitos do CPF; consulta `customers` por `document`; responde `400` (CPF inválido), `404` (não cadastrado), `403` (`status != active`) ou `200` com JWT HS256 de **1 hora** e claims `sub` (id do cliente), `cpf`, `name`, `role = customer`, `iss`.
2. **Segredo único** em `/oficina-api/jwt_secret` (SSM SecureString, gerado pelo Terraform do repositório k8s). A aplicação recebe via Secret do Kubernetes criado pelo CD; a Lambda recebe como variável de ambiente resolvida pelo Terraform.
3. **Lambda `authorizer`** (REQUEST, payload 2.0, resposta simples, cache 300 s pelo header `Authorization`) protege `ANY /api/v1/{proxy+}`. Rotas públicas são declaradas explicitamente no gateway (login, register, status da OS, webhook, health, swagger).
4. **A aplicação continua validando** o JWT no middleware. O gateway é a primeira barreira; a aplicação, a definitiva — permite também chamar o NLB diretamente em testes internos.
5. Não há senha para clientes: o CPF de um cliente **ativo** é a credencial, coerente com o fluxo de atendimento presencial da oficina. A ativação/inativação é controlada pelo operador (`PATCH /customers/{id}/status`).

## Consequências

- **Positivas:** zero alteração no middleware existente; um único formato de token para as duas personas; requisito de "existência e status" mapeado diretamente na coluna `customers.status`; authorizer dá 401 antes de consumir CPU do cluster.
- **Negativas:** o CPF sozinho é uma credencial fraca — aceitável para o escopo acadêmico e mitigado por token curto (1 h), rotas de consulta limitadas e possibilidade de inativar o cliente. HS256 exige que todo verificador conheça o segredo (ver ADR-003). Autorização por papel (`role`) ainda não restringe rotas administrativas para clientes.

## Perguntas abertas

- Restringir rotas administrativas a `role = operator` no authorizer ou no middleware? Próxima evolução natural.
- Adicionar segundo fator (código por e-mail via Resend) para clientes? Fora do escopo desta fase.
