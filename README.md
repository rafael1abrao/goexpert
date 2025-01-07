# Multithreading CEP Lookup

Este projeto em Go demonstra como buscar dados de CEP de forma concorrente em várias APIs diferentes, retornando sempre o resultado da API que responder mais rápido. Caso nenhuma retorne em até **1 segundo**, ocorre um **timeout**.

## Sumário

- [Arquitetura](#arquitetura)
- [Principais Funcionalidades](#principais-funcionalidades)
- [Estrutura de Pastas](#estrutura-de-pastas)
- [Uso](#uso)
- [Configurações e Customizações](#configurações-e-customizações)
- [Testes](#testes)
- [Melhorias Futuras](#melhorias-futuras)

---

## Arquitetura

O projeto segue princípios de **Clean Code** e **Clean Architecture**, além de aplicar o **SOLID**:

1. **Domain (`domain`)**  
   - Contém o modelo de domínio (`CepResponse`) e a lógica principal de busca concorrente (`FetchFasterCep`), incluindo a **validação/normalização de CEP** e controle de timeout (1 segundo).

2. **Providers (`providers`)**  
   - Cada provider implementa a interface `CepProvider` (do pacote `domain`).  
   - Exemplo: `BrasilAPIProvider` e `ViaCepProvider`.  
   - Cada um deles faz o mapeamento do JSON específico de cada API para o modelo de domínio.

3. **Handlers (`handlers`)**  
   - Camada opcional onde ficam as funções que lidam com **requisições HTTP** (ou outro protocolo) recebidas pela aplicação.  
   - Aqui você pode ter, por exemplo, um `cep_handler.go` que implementa um endpoint (`/cep/{cep}`) para receber o CEP, chamar `FetchFasterCep` e retornar o resultado em JSON ao cliente.  
   - A vantagem é manter a regra de negócio no `domain` e apenas orquestrar requisições/respostas na camada de handlers.

4. **Aplicação (`cmd/main.go`)**  
   - Ponto de entrada (função `main`).  
   - Pode apenas rodar um CLI para testes, ou subir um servidor HTTP que expõe as rotas definidas em `handlers`.

Com isso, **o domínio não conhece detalhes de implementação** de cada provider ou de como as requisições são tratadas, e cada provider fica responsável apenas por se comunicar com sua API e converter o retorno em algo que o domínio entenda.

---

## Principais Funcionalidades

- **Busca Concorrente**: Dispara requisições simultâneas a múltiplas APIs.  
- **Retorno da Primeira Resposta**: A aplicação aguarda quem responder primeiro (BrasilAPI, ViaCEP etc.).  
- **Timeout**: Limite de 1 segundo para as requisições completarem. Se exceder, retorna erro de timeout.  
- **Validação de CEP**: Remove hífens e caracteres não numéricos, garantindo que o CEP possua exatamente 8 dígitos antes de consultar as APIs.  
- **Modelagem de Domínio**: As diferentes estruturas de JSON de cada API são convertidas para um modelo comum (`CepResponse`).

---


