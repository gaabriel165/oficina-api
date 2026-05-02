# Linguagem Ubíqua — Sistema de Oficina Mecânica

## Entidades principais

| Termo | Definição |
|---|---|
| **Cliente** | Pessoa física (CPF) ou jurídica (CNPJ) que traz um veículo para atendimento na oficina |
| **Veículo** | Automóvel identificado por placa, marca, modelo e ano, vinculado a um Cliente |
| **Ordem de Serviço (OS)** | Documento central que registra todo o ciclo de atendimento de um Veículo, desde a Recepção até a Entrega |
| **Diagnóstico** | Etapa técnica em que o Mecânico examina o Veículo, identifica os Problemas e define os Serviços e Peças necessários |
| **Orçamento** | Valor total calculado automaticamente a partir dos Serviços e Peças da OS, enviado ao Cliente para Aprovação |
| **Serviço** | Trabalho técnico a ser executado na OS (ex: troca de óleo, alinhamento), com descrição, valor de mão de obra e tempo estimado |
| **Peça** | Item físico utilizado na execução de um Serviço (ex: filtro de óleo, pastilha de freio), controlado pelo Estoque |
| **Estoque** | Controle de quantidade disponível de cada Peça na oficina |

## Atores

| Termo | Definição |
|---|---|
| **Atendente** | Funcionário responsável pela Recepção do Cliente, abertura da OS e Entrega do Veículo |
| **Mecânico** | Profissional técnico responsável pelo Diagnóstico e pela Execução dos Serviços |
| **Administrador** | Usuário com acesso total ao sistema, responsável pela gestão de Clientes, Veículos, Serviços, Peças e OSs |

## Ciclo de vida da OS

| Termo | Definição |
|---|---|
| **Recepção** | Momento de abertura da OS — Cliente identificado, Veículo cadastrado, sintomas registrados |
| **Diagnóstico** | Fase em que o Mecânico examina o Veículo e levanta os Problemas, Serviços e Peças necessários |
| **Problema** | Defeito ou necessidade técnica identificada pelo Mecânico durante o Diagnóstico |
| **Orçamento Gerado** | Estado em que o sistema calculou automaticamente o valor total da OS após o Diagnóstico |
| **Aprovação** | Decisão do Cliente de aceitar ou recusar o Orçamento |
| **Execução** | Fase em que o Mecânico realiza os Serviços aprovados no Veículo |
| **Finalização** | Conclusão da Execução — todos os Serviços foram realizados |
| **Entrega** | Momento em que o Veículo é devolvido ao Cliente após a Finalização |
| **Cancelamento** | Encerramento da OS sem Execução, decorrente da recusa do Orçamento pelo Cliente |

## Status da OS

| Status | Significado |
|---|---|
| **Recebida** | OS aberta, aguardando início do Diagnóstico |
| **Em Diagnóstico** | Mecânico examinando o Veículo |
| **Aguardando Aprovação** | Orçamento enviado ao Cliente, aguardando resposta |
| **Em Execução** | Orçamento aprovado, Serviços em andamento |
| **Finalizada** | Todos os Serviços concluídos, aguardando Entrega |
| **Entregue** | Veículo devolvido ao Cliente |
| **Cancelada** | OS encerrada por recusa do Orçamento |

## Regras de negócio (termos importantes)

| Termo | Definição |
|---|---|
| **Reserva de Peça** | Bloqueio de quantidade no Estoque no momento em que uma Peça é adicionada à OS, evitando conflito com outras OSs |
| **Débito de Estoque** | Baixa definitiva da quantidade reservada no Estoque no início da Execução |
| **Tempo de Execução** | Duração real entre o início e a Finalização da Execução de um Serviço, usado para monitoramento |
| **Valor da OS** | Soma dos valores de Mão de Obra dos Serviços e do custo das Peças utilizadas |
