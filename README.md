# Pokemon API

API para gerenciamento de Pokemons com filtros por tipo.

## Endpoints

### GET /pokemons

Lista pokemons paginados.

Query Parameters:
- `page` (int): Página da lista, default 1
- `limit` (int): Quantidade de itens por página, default 10, máximo 100
- `type` (string): Filtro opcional por tipos de pokemon, valores separados por vírgula. Ex: `fire,water`

Exemplo de requisição:
```
GET /pokemons?type=fire,water&page=1&limit=10
```

Resposta:
```
{
  "data": [
    {
      "id": 10,
      "name": "Charmander",
      "type": ["fire"],
      "level": 12,
      "hp": 35
    },
    {
      "id": 25,
      "name": "Squirtle",
      "type": ["water"],
      "level": 10,
      "hp": 40
    }
  ],
  "page": 1,
  "limit": 10,
  "total_items": 20,
  "total_pages": 2
}

```
