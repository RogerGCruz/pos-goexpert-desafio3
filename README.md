# POS Go Expert - Adição da funcionalidade de listagem de pedidos

## Resumo

Este projeto foi atualizado para incluir a funcionalidade de listagem de pedidos nas três interfaces expostas pela aplicação:

- REST
- GraphQL
- gRPC

Além da implementação da funcionalidade, também foram adicionados testes automatizados para validar a listagem no repositório e no caso de uso.

## O que foi alterado

### Caso de uso e DTOs

Foram criados e ajustados os componentes responsáveis pela listagem de pedidos e pelo compartilhamento dos DTOs entre create e list.

Arquivos envolvidos:

- internal/usecase/list_orders.go
- internal/usecase/create_order.go
- internal/usecase/dto/order_dto.go

### Camada de repositório

A interface do repositório passou a suportar a busca de todos os pedidos, e a implementação concreta ganhou o método responsável por essa consulta no banco.

Arquivos envolvidos:

- internal/entity/interface.go
- internal/infra/database/order_repository.go

### REST

Foi adicionado o endpoint de listagem de pedidos via HTTP.

Arquivos envolvidos:

- internal/infra/web/order_handler.go
- api/list_orders.http

Endpoint:

    GET /orders

### GraphQL

Foi adicionada uma query para listar pedidos.

Arquivos envolvidos:

- internal/infra/graph/schema.graphqls
- internal/infra/graph/schema.resolvers.go
- internal/infra/graph/resolver.go
- internal/infra/graph/generated.go

Query adicionada:

    orders: [Order!]!

### gRPC

Foi adicionado o método de listagem no contrato protobuf e no serviço gRPC.

Arquivos envolvidos:

- internal/infra/grpc/protofiles/order.proto
- internal/infra/grpc/service/order_service.go
- internal/infra/grpc/pb/order.pb.go
- internal/infra/grpc/pb/order_grpc.pb.go

Método adicionado:

    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse)

### Inicialização e injeção de dependências

A inicialização da aplicação foi ajustada para injetar o novo caso de uso de listagem nas camadas web, GraphQL e gRPC.

Arquivos envolvidos:

- cmd/ordersystem/main.go
- cmd/ordersystem/wire.go
- cmd/ordersystem/wire_gen.go

## Testes implementados

Foram adicionados testes automatizados para validar a nova funcionalidade.

### Testes de repositório

Arquivo:

- internal/infra/database/order_repository_test.go

Cobertura:

- persistência de pedidos com Save
- leitura de todos os pedidos com GetAll

### Testes de caso de uso

Arquivo:

- internal/usecase/list_orders_test.go

Cobertura:

- retorno correto da lista de pedidos
- propagação de erro do repositório

### Como rodar os testes

Para executar todos os testes do projeto:

    go test ./...

## Infraestrutura local com Docker Compose

Para desenvolvimento local, o projeto utiliza Docker Compose para subir o MySQL e o RabbitMQ.

Observação importante:

- o projeto não usa Dockerfile para subir MySQL e RabbitMQ
- a infraestrutura local está definida em docker-compose.yaml
- a aplicação Go roda localmente fora de container e se conecta aos serviços expostos em localhost

Serviços definidos no docker-compose.yaml:

### MySQL

- imagem: mysql:5.7
- porta exposta: 3306
- database inicial: orders
- usuário root com senha root
- volume de dados em .docker/mysql

### RabbitMQ

- imagem: rabbitmq:3-management
- porta AMQP: 5672
- painel web: 15672
- usuário: guest
- senha: guest

### Como subir a infraestrutura

Na raiz do projeto, execute:

    docker compose up -d

Para verificar os containers:

    docker compose ps

Para derrubar os containers:

    docker compose down

Para derrubar os containers e remover volumes:

    docker compose down -v

## Configuração da aplicação

A aplicação carrega as configurações a partir de um arquivo .env.

Variáveis esperadas:

- DB_DRIVER
- DB_HOST
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME
- WEB_SERVER_PORT
- GRPC_SERVER_PORT
- GRAPHQL_SERVER_PORT

Exemplo de valores compatíveis com o Docker Compose local:

    DB_DRIVER=mysql
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=root
    DB_PASSWORD=root
    DB_NAME=orders
    WEB_SERVER_PORT=8000
    GRPC_SERVER_PORT=50051
    GRAPHQL_SERVER_PORT=8080

Observação:

- a conexão com RabbitMQ no código está apontando para amqp://guest:guest@localhost:5672/
- se a aplicação também fosse executada em container, localhost deixaria de ser o host correto e seria necessário usar o nome do serviço do compose

## Como rodar a aplicação

Com a infraestrutura disponível e o arquivo .env configurado, execute a partir da raiz do projeto:

    go run ./cmd/ordersystem

Importante:

- não execute apenas o arquivo main.go isoladamente
- o projeto depende de arquivos do mesmo package, incluindo código gerado pelo Wire

## Como testar a criação de pedidos

Arquivo de exemplo:

- api/create_order.http

Exemplo de requisição HTTP:

    POST http://localhost:8000/order HTTP/1.1
    Host: localhost:8000
    Content-Type: application/json

    {
        "id":"124",
        "price": 100.5,
        "tax": 0.5
    }

## Como testar a listagem via REST

Arquivo de exemplo:

- api/list_orders.http

Requisição:

    GET http://localhost:8000/orders HTTP/1.1
    Host: localhost:8000
    Content-Type: application/json

Exemplo usando curl:

    curl http://localhost:8000/orders

Resposta esperada:

    [
      {
        "id": "123",
        "price": 100.5,
        "tax": 0.5,
        "final_price": 101
      }
    ]

Se ainda não houver pedidos cadastrados, o retorno esperado é uma lista vazia.

## Como testar via GraphQL

A aplicação expõe:

- Playground GraphQL em /
- endpoint GraphQL em /query

Com GRAPHQL_SERVER_PORT configurada como 8080, por exemplo, abra no navegador:

    http://localhost:8080/

Query para listar pedidos:

    query {
      orders {
        id
        Price
        Tax
        FinalPrice
      }
    }

Observação importante:

- os campos do schema atual estão definidos como Price, Tax e FinalPrice com inicial maiúscula

Exemplo de resposta:

    {
      "data": {
        "orders": [
          {
            "id": "123",
            "Price": 100.5,
            "Tax": 0.5,
            "FinalPrice": 101
          }
        ]
      }
    }

Também é possível testar pelo terminal com curl:

    curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d "{\"query\":\"query { orders { id Price Tax FinalPrice } }\"}"

## Como testar via gRPC com Evans

O serviço gRPC exposto é o OrderService no package pb.

Método de listagem:

    ListOrders(ListOrdersRequest) returns (ListOrdersResponse)

Como o servidor está com reflection habilitado, o Evans consegue inspecionar os serviços sem precisar carregar manualmente o proto.

Com GRPC_SERVER_PORT configurada como 50051, execute:

    evans -r repl -p 50051

No console do Evans, rode:

    show package
    package pb
    show service
    service OrderService
    call ListOrders
    {}

Observação:

- o request de ListOrders é vazio

Exemplo de resposta:

    {
      "orders": [
        {
          "id": "123",
          "price": 100.5,
          "tax": 0.5,
          "final_price": 101
        }
      ]
    }

### Teste opcional de ponta a ponta no Evans

Você também pode criar um pedido antes de listar.

No Evans:

    call CreateOrder
    {
      "id": "abc123",
      "price": 10,
      "tax": 2
    }

Depois:

    call ListOrders
    {}

## Fluxo sugerido para validação manual

1. Subir MySQL e RabbitMQ com Docker Compose.
2. Configurar o arquivo .env.
3. Iniciar a aplicação com go run ./cmd/ordersystem.
4. Criar um pedido via REST, GraphQL ou gRPC.
5. Consultar a listagem via REST.
6. Consultar a listagem via GraphQL.
7. Consultar a listagem via gRPC com Evans.
8. Confirmar que o mesmo pedido aparece em todas as interfaces.

## Observações finais

- se o projeto não subir, verifique primeiro se o .env está presente na raiz do repositório
- verifique se o MySQL e o RabbitMQ foram iniciados corretamente
- verifique se as portas configuradas no .env correspondem ao ambiente local
