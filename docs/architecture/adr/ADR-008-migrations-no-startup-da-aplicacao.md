# ADR-008 — Migrations executadas no startup da aplicação

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [RFC-002](../rfc/RFC-002-escolha-do-banco-de-dados.md), [database.md](../database.md) |

## Contexto

O schema é versionado com `golang-migrate` em `migrations/` (SQL puro, `up`/`down`). Com o banco agora provisionado por um repositório próprio (`oficina-infra-db`) e a aplicação implantada por outro, era preciso decidir **quem aplica as migrations**: o Terraform do banco, um job dedicado no pipeline, um `Job`/`initContainer` no Kubernetes, ou a própria aplicação ao iniciar.

## Decisão

- As migrations continuam sendo aplicadas **pela aplicação no startup** (`database.RunMigrations` em `cmd/api/main.go`), antes de abrir o servidor HTTP.
- Os arquivos SQL são copiados para a imagem Docker (`COPY migrations ./migrations`).
- `golang-migrate` usa **lock advisory** no PostgreSQL, portanto múltiplas réplicas subindo ao mesmo tempo (2 iniciais + HPA) não aplicam a mesma migration duas vezes; as demais aguardam e seguem.
- O repositório `oficina-infra-db` **não executa DDL**: ele provisiona a instância, o banco lógico inicial e publica a `DATABASE_URL`.
- Regra de ouro mantida: **nunca alterar uma migration aplicada** — a Fase 3 adiciona `000010` para status, normalização e índices.

## Consequências

**Positivas**
- Nenhum passo manual ou job extra no pipeline; um `kubectl rollout` já garante schema compatível com a versão implantada.
- Ambiente local (`docker-compose up`) e produção usam exatamente o mesmo mecanismo.
- Independência entre os repositórios de banco e aplicação: o Terraform do banco não precisa de acesso à rede privada nem de cliente `psql`.

**Negativas**
- Uma migration lenta atrasa o startup e pode disparar o `startupProbe` (hoje tolerante: 30 × 5 s). Migrations grandes exigiriam outro mecanismo.
- A imagem da aplicação carrega privilégios de DDL no banco (usuário master do RDS). Em produção real, separar um usuário de migração de um usuário de runtime.
- `down` migrations não são executadas automaticamente; rollback de schema é manual.

## Alternativas rejeitadas

- **Terraform (`null_resource` com `psql`) no repo do banco** — exigiria acesso de rede à subnet privada a partir do runner e acoplaria o Terraform a arquivos SQL de outro repositório.
- **Job Kubernetes / initContainer separado** — mais "correto" para bancos grandes, mas duplica a imagem/lógica de migração e adiciona um recurso a orquestrar antes do `Deployment`.
- **Etapa dedicada no pipeline do app** — o runner do GitHub não alcança o RDS privado sem bastion ou VPN.
