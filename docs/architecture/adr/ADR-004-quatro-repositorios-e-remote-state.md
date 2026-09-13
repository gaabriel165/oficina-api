# ADR-004 — Quatro repositórios e remote state do Terraform no S3

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## Contexto

Na Fase 2 todo o Terraform (VPC, EKS, RDS, ECR, OIDC) vivia na pasta `infra/` do repositório da aplicação, com estado local no notebook do autor. A Fase 3 exige quatro repositórios independentes — Lambda, infra Kubernetes, infra do banco e aplicação — cada um com CI/CD e deploy automático. Isso cria dependências entre estados: o banco precisa da VPC e do security group dos nós; a Lambda precisa das subnets, do SG de acesso ao banco e dos parâmetros no SSM.

## Decisão

1. **Divisão de responsabilidades**

| Repositório | Contém | Estado remoto |
|---|---|---|
| `oficina-infra-k8s` | VPC, EKS, ECR, OIDC provider + 4 roles, parâmetros `jwt_secret` e `webhook_secret`, instalação de `metrics-server` e `nri-bundle` | `k8s/terraform.tfstate` |
| `oficina-infra-db` | RDS, subnet group, parameter group, SGs `db` e `db-client`, parâmetro `database_url` | `db/terraform.tfstate` |
| `oficina-lambda-auth` | Código Go das Lambdas, IAM de execução, API Gateway, authorizer, integrações | `lambda/terraform.tfstate` |
| `oficina-api` | Código da aplicação, Dockerfile, manifests kustomize, pipelines de build e deploy | — (não usa Terraform) |

2. **Remote state** em um bucket S3 único (`oficina-api-terraform-state-728750563430`), versionado, criptografado e com acesso público bloqueado, criado uma vez fora do Terraform. Lock nativo do backend S3 (`use_lockfile = true`, Terraform ≥ 1.10) — sem tabela DynamoDB.

3. **Acoplamento entre repositórios apenas por `terraform_remote_state`** (db lê k8s; lambda lê k8s e db) e por **parâmetros no SSM** (lidos pelo CD da aplicação). Nenhum repositório edita o estado de outro.

4. **Ordem de criação:** k8s → db → deploy da aplicação (cria o NLB) → lambda. **Destruição** na ordem inversa, precedida de `kubectl delete svc oficina-api -n oficina-api` para que o NLB gerenciado pelo Kubernetes não bloqueie a exclusão da VPC.

5. **Bootstrap:** o primeiro `apply` do repo k8s é feito localmente (ele cria o OIDC provider e o próprio role do seu pipeline); a partir daí o CD assume, usando o mesmo estado no S3.

## Consequências

**Positivas**
- Cada repositório tem ciclo de vida, pipeline e revisão próprios; um `plan` no PR mostra exatamente o impacto daquela camada.
- Estado compartilhado e bloqueado: `apply` do notebook e do pipeline nunca divergem.
- Segredos fluem por SSM, não por outputs copiados manualmente para GitHub Secrets (que era o processo da Fase 2).
- Atende ao requisito de segregação do PDF sem duplicar módulos.

**Negativas**
- Ordem de criação/destruição precisa ser respeitada e documentada; um `destroy` fora de ordem falha por dependência.
- O bucket de estado é um recurso "manual" fora do Terraform (documentado nos READMEs).
- Alterar um output consumido por outro repo exige `apply` em cadeia.

## Alternativas rejeitadas

- **Monorepo com workspaces/pastas** — mais simples de coordenar, mas não cumpre o requisito de quatro repositórios.
- **Terraform Cloud/HCP como backend** — locking e histórico melhores, mas adiciona uma conta/serviço externo; o versionamento do S3 cobre o necessário.
- **Terragrunt** para orquestrar a ordem — resolve a cadeia de `apply`, porém mais uma ferramenta para aprender e documentar em um dia.
