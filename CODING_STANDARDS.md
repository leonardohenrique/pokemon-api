# Padrões e Boas Práticas de Codificação para pokemon-api

Este documento define os padrões de codificação que devem ser seguidos no desenvolvimento do pokemon-api para garantir consistência, legibilidade e qualidade do código.

## Linguagem

O projeto é desenvolvido em Go, seguindo as convenções padrão da linguagem Go.

## Convenções Gerais

- Utilizar nomes descritivos e em inglês para variáveis, funções e pacotes.
- Manter funções e métodos curtos e focados em uma única responsabilidade.
- Comentar trechos não triviais de código para facilitar o entendimento.
- Manter indentação e espaçamento consistentes conforme padrão gofmt.

## Estrutura do Código

- Organizar o código em pacotes adequados, como `handlers`, `middleware` e `store`.
- Separar responsabilidades e evitar código duplicado.

## Controle de Erros

- Tratar erros imediatamente após as chamadas que possam retornar erros.
- Utilizar `log` para registrar erros críticos que impeçam a continuação do programa.

## Testes

- Escrever testes unitários para funções e métodos críticos.
- Garantir que os testes sejam simples e independentes.

## Commit

- Escrever mensagens de commit claras e descritivas.
- Seguir o padrão de mensagens de commit: tipo: descrição breve.

Para dúvidas, consulte a equipe mantenedora do projeto.
