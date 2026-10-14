# API Go

API de teste com frontend de busca de Pokémon, sem dependências externas de código. Requer Go 1.27 ou Docker e acesso à internet para consultar a PokéAPI e carregar imagens.

## Executar localmente

```sh
go run .
```

A porta padrão é `8080`. Configure a variável de ambiente `PORT` para alterá-la.

Abra http://localhost:8080 e pesquise um nome como `pikachu`, `bulbasaur` ou `charizard`. A página mostra imagem, tipos, altura, peso e atributos, com estados de carregamento e erro.

## Executar com Docker

```sh
docker build -t go-api .
docker run --rm -p 8080:8080 --name go-api go-api
```

## Endpoints

| Método | Rota | Resposta |
| --- | --- | --- |
| GET | `/` | Frontend de busca |
| GET | `/health` | `{"status":"ok"}` |
| GET | `/api/pokemon/{name}` | Dados do Pokémon em JSON |

```sh
curl http://localhost:8080/api/pokemon/pikachu
curl http://localhost:8080/health
```

## Validar

```sh
go test ./...
go vet ./...
```

O Dockerfile usa compilação em duas etapas, binário estático e execução com usuário sem privilégios. O servidor possui timeouts HTTP e encerramento gracioso ao receber SIGTERM.

O frontend é incorporado ao binário com `go:embed` e consulta a API Go na mesma origem, sem configuração de CORS. A API consulta `https://pokeapi.co/api/v2/pokemon/{name}/`, valida nomes, limita a resposta a 2 MiB e usa timeout de 8 segundos. Respostas bem-sucedidas ficam em um cache em memória limitado a 128 entradas; ao atingir o limite, o cache é limpo. Não há persistência após reiniciar.

Erros da rota de busca: `400` para nome inválido, `404` para Pokémon inexistente e `502` para falha, timeout ou resposta inválida da PokéAPI.
