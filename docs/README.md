# Pokémon API

## Visão geral
Este repositório contém uma API REST para gerenciamento de pokémons, permitindo operações CRUD (Create, Read, Update, Delete). A aplicação é escrita em Go, utiliza um banco de dados PostgreSQL para persistência, e suporta paginação nas consultas.

## Arquitetura
```mermaid
flowchart TD
    Client -->|HTTP Requests| Server[API HTTP Server]
    Server --> Middleware[Middleware (CORS, Logging)]
    Server --> Handler[PokemonHandler]
    Handler --> Store[PokemonStore (Postgres)]
    Store -->|SQL Queries| Database[(Postgres DB)]

subgraph Server Application
    Middleware --> Handler
end
```

## Como rodar
O projeto pode ser construído e executado usando Go e Docker:

### Sem Docker (local)
1. Certifique-se de ter Go e Postgres instalados.
2. Configure o banco de dados e rode a migration `migrations/001_create_pokemons.sql` para criar a tabela de pokémons.
3. Configure a string de conexão ao banco no código (a confirmar onde exatamente).
4. Execute o servidor com `go run .` ou compile com `go build` e rode o binário.

### Com Docker
1. Execute `docker build -t pokemon-api .`
2. Execute o container com `docker run -p 8080:8080 pokemon-api`

O servidor escutará na porta 8080.

## Configuração
- Banco de dados PostgreSQL é requerido.
- Variáveis de ambiente para conexão a confirmar (não encontrado no código analisado).

## Principais fluxos/endpoints
- `GET /pokemons` - Lista pokémons paginados. Parâmetros query: `page` (padrão 1), `limit` (padrão 10, máximo 100).
- `GET /pokemons/{id}` - Obtém pokémon pelo ID.
- `POST /pokemons` - Cria novo pokémon. Recebe JSON com os dados.
- `PUT /pokemons/{id}` - Atualiza pokémon existente pelo ID.
- `DELETE /pokemons/{id}` - Remove pokémon pelo ID.

Os dados de pokémon contêm: id, name (string), type (array de strings com tipos válidos), level (int 1-100), hp (int), photo_url (string opcional).

## Pontos de atenção
- Validação dos dados é feita usando `github.com/go-playground/validator/v10` com mensagens customizadas.
- Middlewares aplicam CORS liberando todos domínios e logging de requests.
- Pagination é calculada com total de itens e total de páginas no response.
- A configuração do banco e eventuais rotas de entrada principais não foram encontradas (possível entrada main.go e configuração a confirmar).

## Arquivos principais
- `internal/handlers/pokemon_handler.go`: lógica dos endpoints HTTP.
- `internal/store/postgres_store.go`: persistência no Postgres.
- `internal/models/pokemon.go`: modelo e estrutura de dados.
- `internal/middleware/*.go`: middlewares para CORS e Logging.
- `migrations/001_create_pokemons.sql`: script de criação da tabela no banco.
- `Dockerfile`: configuração para container Docker.

---

Documento gerado a partir da análise dos arquivos mencionados. Configuração completa da conexão com o banco e o ponto de entrada do servidor são a confirmar.