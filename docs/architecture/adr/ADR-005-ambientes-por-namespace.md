# ADR-005 — Homologação e produção como namespaces no mesmo cluster

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## Contexto

As regras de proteção do PDF pedem deploy automático das branches de **homologação** e **produção**. A infraestrutura roda em uma conta pessoal, custa cerca de US$ 0,50/h com um cluster e é criada e destruída por sessão. Um segundo cluster EKS (control plane US$ 0,10/h + nós + NAT) dobraria o custo e o tempo de provisionamento (~20 min).

## Decisão

- Na aplicação (`oficina-api`), a branch **`homolog`** faz deploy no namespace **`oficina-api-homolog`** e a branch **`main`** no namespace **`oficina-api`** (produção), ambos no **mesmo cluster**.
- Manifests em **kustomize** com `k8s/base` e overlays `k8s/overlays/homolog` e `k8s/overlays/prod`:
  - `prod`: 2 réplicas + HPA (2..6), `Service` tipo `LoadBalancer` (NLB, consumido pelo API Gateway).
  - `homolog`: 1 réplica, sem HPA, `Service` `ClusterIP` (acesso por `port-forward`), mesmos ConfigMap/Secret com `APP_ENV=homolog`.
- Ambos os ambientes apontam para o **mesmo RDS** nesta fase; o campo `environment` vai nos logs para separar dados no New Relic.
- Repositórios de infraestrutura e Lambda têm um único ambiente: `plan` no pull request, `apply` no merge em `main`.
- A branch `main` de todos os repositórios é protegida (sem commits diretos, pull request obrigatório).

## Consequências

**Positivas**
- Custo e tempo de provisionamento inalterados; pipeline de homologação demonstrável no vídeo.
- Isolamento lógico (RBAC, quotas, nomes) suficiente para validar uma versão antes do merge em `main`.
- Um único agente New Relic observa os dois ambientes, com `namespaceName`/`environment` como dimensão.

**Negativas**
- Isolamento fraco: os dois ambientes compartilham nós, versão do Kubernetes e — nesta fase — o mesmo banco. Um bug em homologação pode alterar dados de produção.
- Sem NLB próprio, homologação não é alcançável pelo API Gateway; testes usam `kubectl port-forward`.

## Alternativas rejeitadas

- **Dois clusters (ou duas contas AWS)** — isolamento real, porém dobra custo e tempo; incompatível com a política de subir a infra apenas em janelas de teste.
- **Terraform workspaces para duplicar RDS/Lambda por ambiente** — aumenta o custo e a matriz de estados sem exigência de negócio.

## Evolução prevista

Criar um banco lógico `oficina_mecanica_homolog` no mesmo RDS (ou instância separada) e um parâmetro SSM `database_url_homolog`; o overlay de homologação passaria a consumir esse segredo, eliminando o compartilhamento de dados.
