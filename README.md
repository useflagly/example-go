# example-go

Exemplo de uso do SDK Go do [UseFlagly](https://useflagly.com.br).

## Pré-requisitos

- Go 1.21+

## Instalação

```bash
go mod tidy
```

## Configuração

```bash
export FLAGLY_API_KEY="sua-api-key-aqui"
```

## Executar

```bash
go run main.go
```

## O que o exemplo demonstra

- Health check da API
- Validar um **Feature Flag** com identificador e contexto
- Validar um **Flow**
- Validar um **Cenário**
- Validar uma **parte de Flow**
- Obter o cache de um flag

## Instalação do SDK

```bash
go get github.com/useflagly/sdk-go@latest
```

### Uso básico

```go
import (
    useflagly "github.com/useflagly/sdk-go"
    "github.com/useflagly/sdk-go/models"
)

client := useflagly.New(useflagly.Options{
    Token: "SUA_API_KEY",
})

result, err := client.ValidateFlag(ctx, "meu-flag", models.ValidateBody{
    Identifier: &userID,
    Context:    map[string]any{"plano": "premium"},
}, "production")
```
