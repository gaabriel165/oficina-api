# ADR-003 — JWT HS256 com segredo compartilhado via SSM Parameter Store

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [RFC-003](../rfc/RFC-003-estrategia-de-autenticacao.md) |

## Contexto

Passam a existir três componentes que precisam **emitir ou verificar** o mesmo JWT: a aplicação (emite para operadores e verifica tudo), a Lambda `auth-cpf` (emite para clientes) e a Lambda `authorizer` (verifica na borda). Desde a Fase 1 o token é assinado com **HS256** e um segredo em variável de ambiente. Era preciso decidir se mantínhamos HS256 (segredo simétrico compartilhado) ou migrávamos para RS256 (par de chaves) e como distribuir o material criptográfico.

## Decisão

- Manter **HS256** como algoritmo único.
- O segredo é **gerado pelo Terraform** do repositório `oficina-infra-k8s` (`random_password`, 48 caracteres) e armazenado **uma única vez** em `/oficina-api/jwt_secret` como `SecureString` no SSM Parameter Store.
- Distribuição:
  - **Aplicação:** o CD lê o parâmetro com `aws ssm get-parameter --with-decryption` e cria o `Secret` do Kubernetes; a aplicação lê `JWT_SECRET` do ambiente.
  - **Lambdas:** o Terraform do repositório `oficina-lambda-auth` lê o parâmetro (`data "aws_ssm_parameter"`) e o injeta como variável de ambiente da função.
- O mesmo padrão vale para `/oficina-api/webhook_secret` e `/oficina-api/database_url`.
- Nenhum segredo é armazenado em GitHub Secrets além das chaves de terceiros (Resend, New Relic).

## Autorização por papel

O token carrega a claim `role` (`operator` no login por e-mail e senha, `customer` na Lambda de CPF; tokens sem a claim são tratados como `operator` por compatibilidade). O middleware `RequireRole` da camada de API restringe as rotas administrativas e as transições operacionais a operadores; clientes acessam apenas a listagem e a consulta das próprias ordens de serviço e a aprovação ou recusa do próprio orçamento, com a checagem de dono feita no caso de uso (`ExecuteForCustomer`). Assim, autenticação e autorização ficam desacopladas: a Lambda prova quem é o cliente, a aplicação decide o que ele pode fazer.

## Consequências

**Positivas**
- Zero mudança no middleware existente e nos testes.
- Fonte única de verdade para o segredo; rotacionar é `terraform apply` no repo k8s seguido de novo deploy da aplicação e da Lambda.
- Sem endpoint JWKS, sem gestão de par de chaves, sem dependência de biblioteca extra na Lambda.
- Verificação simétrica é rápida — relevante no authorizer, que roda em toda requisição não cacheada.

**Negativas**
- Todo verificador possui o material para **emitir** tokens: um comprometimento do authorizer ou da aplicação permite forjar tokens. Mitigação: roles IAM mínimos para leitura do parâmetro, segredo forte, tokens curtos (1 h para clientes).
- Segredo presente na configuração da Lambda (criptografado em repouso com KMS gerenciada, visível para quem tem `lambda:GetFunctionConfiguration`).
- Rotação exige reimplantar os três componentes de forma coordenada; tokens antigos ficam inválidos imediatamente.

## Alternativas rejeitadas

- **RS256 com JWKS publicado pela aplicação e JWT authorizer nativo do API Gateway** — separa quem emite de quem verifica, mas obriga a trocar o middleware, gerar/rotacionar chaves e manter um endpoint JWKS. Ganho de segurança real, custo incompatível com o prazo; fica como evolução.
- **Amazon Cognito como emissor** — tokens RS256, fluxo por CPF sem senha não é nativo (custom auth triggers), e o middleware teria de mudar.
- **AWS Secrets Manager** em vez de SSM — oferece rotação automática, porém cobra por segredo/mês; SSM `SecureString` padrão é gratuito e suficiente.
