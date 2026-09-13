# ADR-007 — Lambda authorizer na borda com validação repetida na aplicação

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [RFC-003](../rfc/RFC-003-estrategia-de-autenticacao.md), [RFC-005](../rfc/RFC-005-api-gateway.md), [ADR-003](./ADR-003-jwt-hs256-segredo-compartilhado.md) |

## Contexto

Com o API Gateway na frente do cluster, surgiu a escolha de **onde** validar o JWT: só na aplicação (gateway como proxy transparente), só no gateway (aplicação confia no cabeçalho), ou em ambos. O HTTP API oferece um authorizer JWT nativo, mas ele exige emissor OIDC com JWKS (RS256), incompatível com o HS256 adotado (ADR-003).

## Decisão

- Implementar um **Lambda authorizer** do tipo **REQUEST**, payload 2.0, com **resposta simples** (`isAuthorized: true|false` + `context`), em Go, no repositório `oficina-lambda-auth`.
- Identity source: `$request.header.Authorization`. **Cache de 300 s** por valor do header — um token válido só invoca a Lambda uma vez a cada 5 minutos.
- Ele valida prefixo `Bearer`, assinatura HS256 com o segredo do SSM e expiração; devolve `sub` e `role` no `context` (visíveis para o backend se necessário).
- Aplicado somente à rota `ANY /api/v1/{proxy+}`; as rotas públicas são declaradas explicitamente e vencem por especificidade.
- **A aplicação continua validando o token** no middleware Gin — a decisão do gateway não substitui a da aplicação.

## Consequências

**Positivas**
- Requisições sem token ou com token inválido recebem `401` **antes** de tocar o NLB/cluster, poupando CPU dos nós e ruído nos logs da aplicação.
- Defesa em profundidade: se alguém alcançar o NLB diretamente (ele é público), a aplicação ainda exige o JWT.
- Cache do authorizer torna o custo por requisição desprezível.
- Ponto único para evoluir autorização por papel (`role`) ou lista de revogação sem tocar na aplicação.

**Negativas**
- Validação duplicada = duas cópias da lógica de verificação (Go na Lambda e Go na aplicação); mantidas idênticas por usarem a mesma biblioteca `golang-jwt/v5`.
- Cache por header significa que um token revogado (ex.: cliente inativado) continua aceito pelo gateway por até 5 min; a aplicação não verifica o status do cliente em cada requisição — a inativação vale para a próxima emissão de token.
- Lista de rotas públicas mantida em dois lugares (gateway e `server.go`).

## Alternativas rejeitadas

- **Somente na aplicação** — gateway vira um proxy caro; não demonstra "controle" no gateway como pede o PDF.
- **Somente no gateway, aplicação confiando em cabeçalhos** — exigiria fechar o NLB (VPC Link + NLB interno) para ser seguro; fora do prazo.
- **JWT authorizer nativo (JWKS)** — depende de RS256/emissor OIDC; rejeitado pelas razões do ADR-003.
