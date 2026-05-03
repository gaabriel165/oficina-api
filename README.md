# Tech Challenge — Oficina Mecânica API

## Sobre o Projeto

Uma oficina mecânica de médio porte enfrentava dificuldades para gerenciar seus atendimentos com anotações manuais e planilhas, gerando erros na priorização, falhas no controle de peças e perda de histórico de clientes.

Este projeto é o MVP do back-end do **Sistema Integrado de Atendimento e Execução de Serviços**, desenvolvido como Tech Challenge Fase 1 da pós-graduação em Arquitetura de Software (FIAP SOAT).

A solução foi construída aplicando **Domain-Driven Design (DDD)** com arquitetura em camadas, focando em:
- Gestão completa do ciclo de vida das ordens de serviço
- Controle de clientes, veículos, peças e serviços
- Fluxo de orçamentos com aprovação/rejeição pelo cliente
- Autenticação JWT para operações administrativas
- Consulta pública de status da OS pelo cliente

## Documentação DDD

- **Event Storming e diagramas:** <!-- INSERIR LINK DO MIRO AQUI -->
- **Linguagem Ubíqua:** [`linguagem-ubiqua.md`](./linguagem-ubiqua.md)
- **Event Storming resumido:** [`event-storming.md`](./event-storming.md)

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

## Por que PostgreSQL?

O PostgreSQL foi escolhido pelos seguintes motivos:

- **Modelo relacional adequado ao domínio:** as entidades do sistema (clientes, veículos, ordens de serviço, peças) possuem relacionamentos bem definidos e consistência transacional crítica — cenário ideal para banco relacional.
- **Integridade referencial:** chaves estrangeiras garantem que uma OS não pode referenciar um veículo ou cliente inexistente sem tratamento explícito.
- **Suporte a tipos avançados:** uso de `UUID` como chave primária, `TIMESTAMPTZ` para datas e `NUMERIC` para valores monetários sem perda de precisão.
- **Migrations versionadas:** o ecossistema Go tem suporte maduro para PostgreSQL com `golang-migrate`, permitindo evoluir o schema de forma controlada.
- **Maturidade e confiabilidade:** amplamente adotado em sistemas de produção, com excelente suporte a transações ACID — essencial para operações como aprovação de orçamento e débito de estoque que devem ser atômicas.

## Arquitetura

O projeto segue arquitetura em camadas baseada em DDD:

```
cmd/                        → entrypoint da aplicação
configs/                    → carregamento de variáveis de ambiente
internal/
  domain/
    entity/                 → entidades de domínio com comportamento e regras de negócio
    repository/             → interfaces de repositório e erros de domínio
    valueobject/            → value objects (CPF, CNPJ, Plate, OrderStatus)
  application/
    usecase/                → casos de uso (um por operação)
    mocks/                  → mocks Testify para testes unitários
  infrastructure/
    database/               → conexão com o banco, migrations, models GORM
    repository/             → implementações GORM dos repositórios
    external/               → clientes de APIs externas (Brasil API para validação de CNPJ)
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
A documentação Swagger estará disponível em `http://localhost:8080/swagger/index.html`.

### Variáveis de Ambiente

| Variável | Descrição | Exemplo |
|---|---|---|
| `SERVER_PORT` | Porta do servidor HTTP | `8080` |
| `DATABASE_URL` | String de conexão com o PostgreSQL | `postgres://postgres:postgres@postgres:5432/oficina_mecanica?sslmode=disable` |
| `JWT_SECRET` | Chave secreta para assinatura do JWT | `sua_chave_secreta` |

## Endpoints da API

### Autenticação

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/auth/register` | Não | Criar conta de usuário |
| POST | `/api/v1/auth/login` | Não | Autenticar e receber JWT |

Todas as demais rotas exigem o header:
```
Authorization: Bearer <token>
```

### Clientes

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| GET | `/api/v1/customers` | Sim | Listar todos os clientes |
| POST | `/api/v1/customers` | Sim | Cadastrar cliente (CPF ou CNPJ) |
| GET | `/api/v1/customers/:id` | Sim | Buscar cliente por ID |
| PUT | `/api/v1/customers/:id` | Sim | Atualizar cliente |
| DELETE | `/api/v1/customers/:id` | Sim | Remover cliente |

### Veículos

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| GET | `/api/v1/vehicles` | Sim | Listar todos os veículos |
| POST | `/api/v1/vehicles` | Sim | Cadastrar veículo |
| GET | `/api/v1/vehicles/:id` | Sim | Buscar veículo por ID |
| PUT | `/api/v1/vehicles/:id` | Sim | Atualizar veículo |
| DELETE | `/api/v1/vehicles/:id` | Sim | Remover veículo |

### Peças

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| GET | `/api/v1/parts` | Sim | Listar todas as peças |
| POST | `/api/v1/parts` | Sim | Cadastrar peça |
| GET | `/api/v1/parts/:id` | Sim | Buscar peça por ID |
| PUT | `/api/v1/parts/:id` | Sim | Atualizar peça |
| DELETE | `/api/v1/parts/:id` | Sim | Remover peça |
| PATCH | `/api/v1/parts/:id/stock` | Sim | Adicionar ao estoque |

### Serviços

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| GET | `/api/v1/services` | Sim | Listar todos os serviços |
| POST | `/api/v1/services` | Sim | Cadastrar serviço |
| GET | `/api/v1/services/:id` | Sim | Buscar serviço por ID |
| PUT | `/api/v1/services/:id` | Sim | Atualizar serviço |
| DELETE | `/api/v1/services/:id` | Sim | Remover serviço |

### Ordens de Serviço

| Método | Caminho | Auth | Descrição |
|---|---|---|---|
| GET | `/api/v1/service-orders/:id/status` | **Não** | Consulta pública de status pelo cliente |
| GET | `/api/v1/service-orders` | Sim | Listar todas as ordens de serviço |
| POST | `/api/v1/service-orders` | Sim | Criar ordem de serviço |
| GET | `/api/v1/service-orders/:id` | Sim | Buscar ordem de serviço por ID |
| GET | `/api/v1/service-orders/metrics` | Sim | Tempo médio de execução por serviço |
| PATCH | `/api/v1/service-orders/:id/start-diagnosis` | Sim | Iniciar diagnóstico |
| POST | `/api/v1/service-orders/:id/services` | Sim | Adicionar serviço à ordem |
| POST | `/api/v1/service-orders/:id/parts` | Sim | Adicionar peça à ordem |
| PATCH | `/api/v1/service-orders/:id/send-budget` | Sim | Enviar orçamento ao cliente |
| PATCH | `/api/v1/service-orders/:id/approve-budget` | Sim | Aprovar orçamento |
| PATCH | `/api/v1/service-orders/:id/reject-budget` | Sim | Rejeitar orçamento |
| PATCH | `/api/v1/service-orders/:id/finish-execution` | Sim | Marcar execução como concluída |
| PATCH | `/api/v1/service-orders/:id/deliver` | Sim | Entregar veículo ao cliente |

## Testes

### Testes Unitários

```bash
go test ./internal/domain/... ./internal/application/usecase/...
```

### Testes de Integração

Requerem o Docker em execução. Cada teste sobe um container PostgreSQL isolado via testcontainers-go.

```bash
go test ./internal/integration/... -timeout 120s
```

### Todos os Testes com Cobertura

```bash
go test ./internal/domain/... ./internal/application/usecase/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Desenvolvimento

### Regenerar Documentação Swagger

```bash
~/go/bin/darwin_amd64/swag init -g cmd/api/main.go --output docs
```

### Criar Nova Migration

```bash
migrate create -ext sql -dir migrations -seq <nome_da_migration>
```
