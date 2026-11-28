# Arquitetura do pokemon-api

Este documento descreve a arquitetura do projeto pokemon-api, detalhando seus componentes principais e decisões arquiteturais.

## Visão Geral

O pokemon-api é uma API REST desenvolvida em Go, que oferece funcionalidades para gerenciamento de registros de pokemons.

## Componentes

- **main.go**: Ponto de entrada da aplicação, configura rotas HTTP e middlewares.
- **handlers**: Contém os manipuladores das requisições HTTP relacionados aos recursos de pokemon.
- **middleware**: Implementa funcionalidades transversais para as requisições, como logging e CORS.
- **store**: Responsável pela interface com o banco de dados PostgreSQL, usando pooling para conexões.

## Fluxo de Requisição

1. O cliente faz uma requisição HTTP para a API.
2. O middleware executa operações como logging e controle de acesso CORS.
3. O handler correspondente processa a requisição, interage com a camada de store para acessar dados.
4. O handler retorna a resposta ao cliente.

## Banco de Dados

A aplicação conecta-se a um banco PostgreSQL configurado via variável de ambiente DATABASE_URL.

## Decisões Arquiteturais

- Uso de Go para alta performance e simplicidade.
- Separação clara de responsabilidades para facilitar manutenção e testes.
- Encadeamento de middlewares para processamento modular das requisições.
- Exposição de endpoints RESTful para operações CRUD.
