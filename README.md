# API Go

API de teste sem dependências externas. Requer Go 1.27 ou Docker.

## Executar localmente

```sh
go run .
```

A porta padrão é `8080`. Configure a variável de ambiente `PORT` para alterá-la.

## Executar com Docker

```sh
docker build -t go-api .
docker run --rm -p 8080:8080 --name go-api go-api
```

## Endpoints

| Método | Rota | Resposta |
| --- | --- | --- |
| GET | `/` | `{"message":"API Go funcionando"}` |
| GET | `/health` | `{"status":"ok"}` |

```sh
curl http://localhost:8080/
curl http://localhost:8080/health
```

## Validar

```sh
go test ./...
go vet ./...
```

O Dockerfile usa compilação em duas etapas, binário estático e execução com usuário sem privilégios. O servidor possui timeouts HTTP e encerramento gracioso ao receber SIGTERM.
