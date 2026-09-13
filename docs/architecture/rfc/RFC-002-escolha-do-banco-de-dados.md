# RFC-002 — Escolha do banco de dados

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |
| **Relacionado** | [database.md](../database.md) |

## Resumo

Manter **PostgreSQL 16** como banco único do sistema, gerenciado pelo **Amazon RDS**, provisionado por Terraform em repositório próprio (`oficina-infra-db`), com a string de conexão publicada no SSM Parameter Store.

## Contexto / Problema

O domínio é transacional: uma ordem de serviço agrega itens e peças, referencia cliente e veículo, e sua aprovação debita estoque. A Fase 3 pede um banco gerenciado, documentação e justificativa formal do modelo, além de ajustes de consistência e performance. Também surge um segundo consumidor do banco: a Lambda de autenticação, que precisa consultar clientes fora do cluster.

## Opções consideradas

| Critério | PostgreSQL (RDS) | MySQL (RDS) | Aurora Serverless v2 (PG) | DynamoDB | MongoDB Atlas |
|---|---|---|---|---|---|
| Transações multi-tabela e FKs | Sim | Sim | Sim | Não (transações limitadas, sem FK) | Parcial |
| Tipos exatos (`NUMERIC`, `UUID`) | Nativos | `DECIMAL`; UUID sem tipo nativo | Nativos | Não | Decimal128 |
| Consultas analíticas (média de tempo por serviço, ordenação por prioridade) | SQL completo | SQL completo | SQL completo | Na aplicação | Aggregation |
| Custo em uso intermitente | `db.t3.micro` ≈ US$ 0,017/h | igual | mínimo 0,5 ACU ≈ US$ 0,06/h | quase zero | fora da VPC |
| Tempo de criação/destruição | ~8 min | ~8 min | ~10 min | segundos | n/a |
| Acesso pela Lambda na VPC | SG + `pgx` | SG + driver | SG + `pgx` | SDK, sem VPC | Internet |
| Aderência ao código existente (GORM, migrations) | **Total** | Alta (ajustar dialeto) | Total | Reescrita | Reescrita |
| Risco | Baixo | Baixo | Baixo/médio | Alto | Alto |

## Proposta

1. **PostgreSQL 16.15 em RDS**, `db.t3.micro`, gp3 20 GB criptografado, single-AZ, `publicly_accessible = false`.
2. Security group do banco aceita 5432 apenas do SG dos nós do EKS e de um SG `db-client` anexado à Lambda.
3. Observabilidade do banco: Performance Insights, export de logs para o CloudWatch e `log_min_duration_statement = 500`.
4. `DATABASE_URL` completa (com `sslmode=require`) publicada como `SecureString` em `/oficina-api/database_url`; lida pelo CD da aplicação e pelo Terraform da Lambda — nenhuma senha em GitHub Secrets.
5. Ajustes de modelo na migration `000010`: coluna `customers.status`, normalização de `document_type` e índices nas colunas de FK e de status (detalhes em [database.md](../database.md)).
6. Backup automático de 1 dia e `skip_final_snapshot = true`, coerentes com um ambiente que é destruído ao final de cada sessão.

## Consequências

- **Positivas:** integridade garantida pelo banco; transação de aprovação de orçamento sem código compensatório; ajustes de índice atendem exatamente às consultas da aplicação; segredo do banco fora do GitHub; Lambda e app compartilham o mesmo dado de cliente sem replicação.
- **Negativas:** single-AZ e retenção curta não servem para produção real (mitigação: são variáveis Terraform — `multi_az`, `db_backup_retention_days`); dois consumidores do mesmo schema acoplam a Lambda ao modelo de `customers` (mitigação: a Lambda lê apenas `id`, `name`, `status` por `document`, uma consulta estável e indexada).
- **Migrations continuam no startup da aplicação** (ADR-008); o repositório de infra do banco não executa DDL.

## Perguntas abertas

- Separar um usuário de banco somente-leitura para a Lambda? Desejável em produção; adiado por prazo.
- Um banco lógico separado para homologação? Hoje ambos os namespaces apontam para o mesmo banco (ver ADR-005).
