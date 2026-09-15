# portfolio-backend

Monorepo proyek portfolio backend engineering. Setiap service di `services/` dibangun bertahap sebagai bagian dari roadmap: clean architecture, testing, event-driven design, observability, sampai deployment production-grade.

## Struktur Repo
├── cmd/portfolio/ # skeleton/smoke-test awal (Fase 0)
├── services/ # tiap service independen (diisi mulai Fase 1)
├── docs/adr/ # Architecture Decision Records
├── .golangci.yml # konfigurasi lint
└── .github/workflows/ # CI pipeline

## Prasyarat

- Go 1.23+
- golangci-lint v2

## Cara Menjalankan

```bash
go mod tidy
go run ./cmd/portfolio
```

## Testing & Lint

```bash
golangci-lint run ./...
go test ./... -race -cover
```

## Services

| Service | Status | Deskripsi |
|---|---|---|
| url-shortener | Done | Clean architecture, unit+integration test, REST API |

## Architecture Decision Records

Keputusan desain penting didokumentasikan di [`docs/adr/`](./docs/adr/).
