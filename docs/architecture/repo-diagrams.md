# Diagramas por repositório

| | |
|---|---|
| **Status** | Vigente |
| **Data** | 2026-09-13 |
| **Autor** | Gabriel Camargo |

Cada bloco abaixo é o diagrama de arquitetura específico de um repositório, pronto para ser colado na seção "Arquitetura" do respectivo `README.md`.

## 1. `oficina-api` — aplicação principal

```mermaid
flowchart LR
    subgraph GitHub["GitHub - oficina-api"]
        PR["Pull Request"] --> CI["CI: build, vet, testes unitários e integração"]
        Homolog["branch homolog"] --> CDh["CD homolog"]
        Main["branch main"] --> CDp["CD prod"]
    end

    subgraph AWS["AWS"]
        ECR[("ECR")]
        SSM[("SSM: database_url, jwt_secret, webhook_secret")]
        subgraph EKS["EKS"]
            subgraph NsH["namespace oficina-api-homolog"]
                DepH["Deployment 1 réplica"]
            end
            subgraph NsP["namespace oficina-api"]
                DepP["Deployment 2..6 réplicas + HPA"]
                Svc["Service LoadBalancer - NLB"]
            end
        end
        RDS[("RDS PostgreSQL")]
    end

    NR["New Relic - traces OTLP e logs JSON"]
    Resend["Resend"]
    BrasilAPI["BrasilAPI"]

    CDh -->|"docker build e push"| ECR
    CDp -->|"docker build e push"| ECR
    CDh -.->|"lê segredos"| SSM
    CDp -.->|"lê segredos"| SSM
    CDh -->|"kubectl apply -k overlays/homolog"| DepH
    CDp -->|"kubectl apply -k overlays/prod"| DepP
    Svc --> DepP
    DepP --> RDS
    DepH --> RDS
    DepP --> Resend
    DepP --> BrasilAPI
    DepP -.-> NR
    DepH -.-> NR
```

## 2. `oficina-infra-k8s` — infraestrutura Kubernetes

```mermaid
flowchart LR
    subgraph GitHub["GitHub - oficina-infra-k8s"]
        PR["Pull Request"] --> Plan["terraform fmt, validate, plan"]
        Main["branch main"] --> Apply["terraform apply + helm"]
    end

    S3[("S3 state - k8s/")]

    subgraph AWS["AWS us-east-1"]
        subgraph VPC["VPC 10.0.0.0/16"]
            Pub["Subnets públicas - NAT Gateway"]
            Priv["Subnets privadas"]
            subgraph EKS["EKS 1.34"]
                NG["Managed node group t3.small 2..4 - AL2023"]
                MS["metrics-server"]
                NRB["New Relic nri-bundle"]
            end
        end
        ECR[("ECR oficina-api")]
        OIDC["OIDC provider GitHub"]
        Roles["4 roles IAM: app, k8s-infra, db-infra, lambda"]
        SSM[("SSM: jwt_secret, webhook_secret")]
    end

    Apply -.-> S3
    Apply --> VPC
    Apply --> EKS
    Apply --> ECR
    Apply --> OIDC
    Apply --> Roles
    Apply --> SSM
    Apply -->|"helm upgrade --install"| MS
    Apply -->|"helm upgrade --install"| NRB
    OIDC --> Roles
    Roles -.->|"assumidas via OIDC"| GitHub
    S3 -.->|"remote state lido por"| Consumers["oficina-infra-db e oficina-lambda-auth"]
```

## 3. `oficina-infra-db` — banco de dados gerenciado

```mermaid
flowchart LR
    subgraph GitHub["GitHub - oficina-infra-db"]
        PR["Pull Request"] --> Plan["terraform fmt, validate, plan"]
        Main["branch main"] --> Apply["terraform apply"]
    end

    S3k[("S3 state - k8s/")]
    S3d[("S3 state - db/")]

    subgraph AWS["AWS us-east-1"]
        subgraph VPC["VPC - subnets privadas"]
            RDS[("RDS PostgreSQL 16.15 db.t3.micro<br/>gp3 20GB criptografado, backup 1 dia")]
            SGdb["SG db: 5432 apenas de SG nós EKS e SG db-client"]
            SGclient["SG db-client - anexado à Lambda"]
            PG["Parameter group: log_min_duration_statement 500"]
        end
        SSM[("SSM: /oficina-api/database_url")]
        CW["CloudWatch: logs postgresql, Performance Insights"]
    end

    Apply -.->|"lê vpc, subnets, SG dos nós"| S3k
    Apply -.-> S3d
    Apply --> RDS
    Apply --> SGdb
    Apply --> SGclient
    Apply --> PG
    Apply --> SSM
    RDS -.-> CW
    SSM -.->|"lido pelo CD do app e pela Lambda"| Consumers["oficina-api e oficina-lambda-auth"]
```

## 4. `oficina-lambda-auth` — function serverless e API Gateway

```mermaid
flowchart LR
    subgraph GitHub["GitHub - oficina-lambda-auth"]
        PR["Pull Request"] --> CI["go test, go build arm64, terraform plan"]
        Main["branch main"] --> CD["build zips + terraform apply"]
    end

    S3["S3 state - k8s/, db/, lambda/"]

    subgraph AWS["AWS us-east-1"]
        APIGW["API Gateway HTTP API"]
        subgraph VPC["VPC - subnets privadas"]
            LA["Lambda auth-cpf<br/>provided.al2023 arm64"]
            RDS[("RDS PostgreSQL")]
        end
        LZ["Lambda authorizer<br/>REQUEST, simple response, cache 300s"]
        SSM[("SSM: jwt_secret, database_url")]
        NLB["NLB do Service Kubernetes<br/>descoberto por tag"]
        CW["CloudWatch Logs"]
    end

    Cliente["Cliente"] -->|"POST /auth/cpf"| APIGW
    Operador["Operador"] -->|"/api/v1/** Bearer JWT"| APIGW
    APIGW -->|"AWS_PROXY"| LA
    APIGW -.->|"ANY /api/v1/{proxy+}"| LZ
    APIGW -->|"HTTP_PROXY"| NLB
    LA --> RDS
    CD -.-> S3
    CD --> LA
    CD --> LZ
    CD --> APIGW
    LA -.->|"env resolvido pelo Terraform"| SSM
    LZ -.->|"env resolvido pelo Terraform"| SSM
    LA -.-> CW
    LZ -.-> CW
```
