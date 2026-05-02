# Tech Challenge — Oficina Mecânica API

API REST para gerenciamento de uma oficina mecânica, desenvolvida como Tech Challenge Fase 1 da FIAP SOAT.

## Stack

- **Go 1.25** — linguagem principal
- **Gin** — framework HTTP
- **GORM** — ORM para PostgreSQL
- **PostgreSQL 16** — banco de dados relacional
- **golang-migrate** — migrations versionadas
- **golang-jwt/jwt** — autenticação via JWT
- **swaggo/swag** — documentação Swagger
- **Testify** — assertions para testes unitários
- **Mockery** — geração automática de mocks a partir das interfaces de repositório
- **testcontainers-go** — testes de integração com PostgreSQL real

## Arquitetura

O projeto segue arquitetura em camadas baseada em DDD:

```
cmd/                        → entrypoint da aplicação
configs/                    → carregamento de variáveis de ambiente
internal/
  domain/
    entity/                 → entidades de domínio com comportamento e regras de negócio
    repository/             → interfaces de repositório e erros de domínio
    valueobject/            → value objects (ex: OrderStatus)
  application/
    usecase/                → casos de uso (um por operação)
    mocks/                  → mocks Testify para testes unitários
  infrastructure/
    database/               → conexão com o banco, migrations, models GORM
    repository/             → implementações GORM dos repositórios
    external/               → clientes de APIs externas (ex: Brasil API para CNPJ)
  api/
    handler/                → handlers Gin (rotas definidas por domínio)
    middleware/             → middleware JWT de autenticação
    dto/                    → DTOs de request/response
    response/               → mapeador centralizado de erros para HTTP
    server.go               → montagem do servidor (rotas, repositórios, casos de uso)
migrations/                 → migrations SQL versionadas (golang-migrate)
docs/                       → documentação Swagger (gerada automaticamente)
```

## Conceitos de Domínio

| Entidade | Descrição |
|---|---|
| **User** | Operador do sistema — autenticado via e-mail/senha, recebe JWT |
| **Customer** | Proprietário do veículo — identificado por CPF ou CNPJ |
| **Vehicle** | Pertence a um Customer — identificado pela placa (formato antigo ou Mercosul) |
| **Service** | Mão de obra oferecida pela oficina (ex: Troca de Óleo) |
| **Part** | Peça física com controle de estoque (ex: Filtro de Óleo) |
| **ServiceOrder** | Agregado principal — controla todo o ciclo de vida de um reparo |

### Fluxo de Status da Ordem de Serviço

```
received → in_diagnosis → waiting_approval → in_execution → finished → delivered
                                          ↘ cancelled (quando orçamento é rejeitado)
```

## Como Executar

### Pré-requisitos

- Docker e Docker Compose

### Inicialização

```bash
# 1. Copiar e configurar as variáveis de ambiente
cp .env.example .env

# 2. Subir a aplicação e o banco de dados
docker-compose up
```

A API estará disponível em `http://localhost:8080`.

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|---|---|---|
| `SERVER_PORT` | Porta do servidor HTTP | `8080` |
| `DATABASE_URL` | String de conexão com o PostgreSQL | — |
| `JWT_SECRET` | Chave secreta para assinatura do JWT | — |

## Documentação da API

A interface Swagger estará disponível em `http://localhost:8080/swagger/index.html` após subir a aplicação.

### Autenticação

Todas as rotas, exceto `POST /api/v1/auth/register` e `POST /api/v1/auth/login`, exigem um Bearer token:

```
Authorization: Bearer <token>
```

### Resumo dos Endpoints

| Método | Caminho | Descrição |
|---|---|---|
| POST | `/api/v1/auth/register` | Criar conta de usuário |
| POST | `/api/v1/auth/login` | Autenticar e receber JWT |
| GET | `/api/v1/customers` | Listar todos os clientes |
| POST | `/api/v1/customers` | Cadastrar cliente (CPF ou CNPJ) |
| GET | `/api/v1/customers/:id` | Buscar cliente por ID |
| PUT | `/api/v1/customers/:id` | Atualizar cliente |
| DELETE | `/api/v1/customers/:id` | Remover cliente |
| GET | `/api/v1/vehicles` | Listar todos os veículos |
| POST | `/api/v1/vehicles` | Cadastrar veículo |
| GET | `/api/v1/vehicles/:id` | Buscar veículo por ID |
| PUT | `/api/v1/vehicles/:id` | Atualizar veículo |
| DELETE | `/api/v1/vehicles/:id` | Remover veículo |
| GET | `/api/v1/parts` | Listar todas as peças |
| POST | `/api/v1/parts` | Cadastrar peça |
| GET | `/api/v1/parts/:id` | Buscar peça por ID |
| PUT | `/api/v1/parts/:id` | Atualizar peça |
| DELETE | `/api/v1/parts/:id` | Remover peça |
| PATCH | `/api/v1/parts/:id/stock` | Adicionar ao estoque |
| GET | `/api/v1/services` | Listar todos os serviços |
| POST | `/api/v1/services` | Cadastrar serviço |
| GET | `/api/v1/services/:id` | Buscar serviço por ID |
| PUT | `/api/v1/services/:id` | Atualizar serviço |
| DELETE | `/api/v1/services/:id` | Remover serviço |
| GET | `/api/v1/service-orders` | Listar todas as ordens de serviço |
| POST | `/api/v1/service-orders` | Criar ordem de serviço |
| GET | `/api/v1/service-orders/:id` | Buscar ordem de serviço por ID |
| POST | `/api/v1/service-orders/:id/start-diagnosis` | Iniciar diagnóstico |
| POST | `/api/v1/service-orders/:id/services` | Adicionar serviço à ordem |
| POST | `/api/v1/service-orders/:id/parts` | Adicionar peça à ordem |
| POST | `/api/v1/service-orders/:id/send-budget` | Enviar orçamento ao cliente |
| POST | `/api/v1/service-orders/:id/approve-budget` | Cliente aprova o orçamento |
| POST | `/api/v1/service-orders/:id/reject-budget` | Cliente rejeita o orçamento |
| POST | `/api/v1/service-orders/:id/finish-execution` | Marcar execução como concluída |
| POST | `/api/v1/service-orders/:id/deliver` | Entregar veículo ao cliente |

## Testes

### Testes Unitários

```bash
go test ./internal/application/usecase/... ./internal/domain/entity/...
```

### Testes de Integração

Requerem o Docker em execução. Cada teste sobe um container PostgreSQL isolado via testcontainers-go.

```bash
go test ./internal/integration/... -timeout 120s
```

### Todos os Testes

```bash
go test ./...
```

## Desenvolvimento

### Regenerar Documentação Swagger

```bash
~/go/bin/swag init -g cmd/main.go
```

### Criar Nova Migration

```bash
migrate create -ext sql -dir migrations -seq <nome_da_migration>
```

