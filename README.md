# example-go

Exemplo de uso do SDK Go do [UseFlagly](https://useflagly.com.br).

## Pré-requisitos

- Go 1.21+ (ou Docker)

## Configuração

Crie um arquivo `.env` (ou copie o `.env` da raiz do repositório de exemplos):

```env
FLAGLY_API_KEY=sua-api-key-aqui
FLAGLY_IDENTIFIER=seu-identifier
FLAGLY_SLUG=seu-slug
FLAGLY_ENVIRONMENT=HML
```

## Executar com Docker

```bash
docker build -t example-go .
docker run --rm --env-file .env example-go
```

## Executar localmente

```bash
go mod tidy
go run main.go
```

## O que o exemplo demonstra

1. **Health check** da API
2. **Initialize** — registra o identifier+slug e inicia a avaliação assíncrona
3. **GetResult** — obtém a árvore de resultados com todos os slugs avaliados
4. Itera o resultado chamando **ValidateFlow**, **ValidateFlowPart** e **ValidateFlag** com os slugs reais

## SDK

```bash
go get github.com/useflagly/sdk-go@latest
```

```go
import (
    useflagly "github.com/useflagly/sdk-go"
    "github.com/useflagly/sdk-go/models"
)

client := useflagly.New(useflagly.Options{Token: "SUA_API_KEY"})

// 1. Inicializar
client.Initialize(ctx, models.ReceiveMessage{
    Identifier: "user-123",
    Slug:       "meu-slug",
}, "HML")

// 2. Obter resultado
result, _ := client.GetResult(ctx, "user-123")

// 3. Validar flags
client.ValidateFlag(ctx, "meu-flag", models.ValidateBody{
    Identifier: &userID,
}, "HML")
```
