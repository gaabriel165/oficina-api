# RFC-001 — Escolha da nuvem

| | |
|---|---|
| **Status** | Aceita |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

## Resumo

Manter a **AWS (região `us-east-1`)** como provedor único para Kubernetes gerenciado, banco gerenciado, function serverless, API Gateway e registry de imagens.

## Contexto / Problema

A Fase 3 exige API Gateway, function serverless, banco gerenciado, cluster Kubernetes com escalabilidade e Terraform, em quatro repositórios com deploy automático. A Fase 2 já provisionou EKS, RDS e ECR na AWS com Terraform e pipelines OIDC funcionando. O ambiente é pago do bolso do autor e é criado e destruído por sessão de teste; qualquer troca de provedor implicaria reescrever infra validada e reaprender operação em um dia de prazo.

## Opções consideradas

| Critério | AWS | GCP | Azure |
|---|---|---|---|
| Aderência aos requisitos (K8s, DB gerenciado, serverless, gateway, IaC) | EKS, RDS, Lambda, API Gateway, Terraform — todos maduros | GKE, Cloud SQL, Cloud Functions/Run, API Gateway | AKS, Azure DB, Functions, APIM |
| Esforço de migração | **Zero** — infra da Fase 2 reaproveitada | Alto — reescrever VPC/cluster/DB/OIDC | Alto |
| Custo em uso intermitente (sobe/desce por sessão) | Control plane EKS US$ 0,10/h + 2× t3.small + db.t3.micro + NAT ≈ US$ 0,50/h | GKE Autopilot cobra por pod; cluster zonal padrão tem 1 gratuito por conta | AKS control plane gratuito, mas APIM caro no tier básico |
| Serverless com acesso a banco privado na VPC | Lambda em VPC + SG — padrão bem documentado | Cloud Functions precisa de Serverless VPC Connector | Functions precisa de VNet integration (plano Premium) |
| Autenticação de pipeline sem chave estática | OIDC GitHub → IAM role já provisionado | Workload Identity Federation | Federated credentials |
| Risco operacional para o prazo | Baixo — comandos e falhas já conhecidos | Médio/alto | Médio/alto |

### Região

`us-east-1` é a região mais barata da AWS, é onde todo o Terraform e as variáveis do GitHub já apontam e não há usuários reais sensíveis a latência. `sa-east-1` custaria cerca de 40% a mais pelos mesmos recursos sem benefício para uma demonstração.

## Proposta

1. Manter AWS `us-east-1` como único provedor.
2. Segregar a infra em `oficina-infra-k8s` (rede, EKS, ECR, IAM/OIDC, SSM) e `oficina-infra-db` (RDS), ambos em Terraform com estado remoto no S3 (ver ADR-004).
3. Implementar API Gateway (HTTP API) e Lambda em `oficina-lambda-auth`, também em Terraform.
4. Atualizar o EKS para **1.34**: a versão 1.31 da Fase 2 entrou em *extended support*, cujo control plane custa seis vezes mais.

## Consequências

- **Positivas:** reaproveitamento total da Fase 2; um único mecanismo de identidade (IAM/OIDC) para os quatro pipelines; Lambda e RDS na mesma VPC sem conectores adicionais; documentação e operação já conhecidas.
- **Negativas:** lock-in em serviços AWS (API Gateway, Lambda, SSM). Mitigação: o código da aplicação não conhece a nuvem — a Lambda é Go puro atrás do `aws-lambda-go`, e a aplicação lê apenas variáveis de ambiente.
- **Custo:** ambiente completo ≈ US$ 0,50/h; a política é destruir tudo fora das janelas de teste.

## Perguntas abertas

- Vale mover o estado do Terraform para um backend com histórico (Terraform Cloud) se o time crescer? Hoje o versionamento do bucket S3 é suficiente.
