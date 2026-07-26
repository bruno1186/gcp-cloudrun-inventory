# gcp-cloudrun-inventory

Microservico de inventario em **Go**, projetado para rodar no **Google Cloud Run**, containerizado com Docker (imagem distroless) e publicado no GitHub Container Registry.

## Endpoints

- GET /health - healthcheck do servico
- GET /items - lista os itens do inventario
- POST /items - cria ou atualiza um item (sku, name, quantity)

## Rodando localmente

```bash
go run .
```

```bash
go test ./... -v
```

## Docker

```bash
docker build -t gcp-cloudrun-inventory .
docker run -p 8080:8080 gcp-cloudrun-inventory
```

## Deploy no Cloud Run

```bash
gcloud run deploy gcp-cloudrun-inventory \
  --image ghcr.io/bruno1186/gcp-cloudrun-inventory:latest \
  --platform managed \
  --region southamerica-east1
```

## CI/CD

O workflow builda e testa a aplicacao Go e, em pushes para main, publica a imagem Docker no GHCR (ghcr.io/bruno1186/gcp-cloudrun-inventory).

## Stack

Go, Google Cloud Run, Docker, GitHub Actions, GHCR
