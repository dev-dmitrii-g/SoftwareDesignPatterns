# Hospital Admission

A simple Go and SQLite website that uses a facade to coordinate insurance verification, equipment reservation, and pharmacy dispatch. If a later step fails, the facade reverses completed steps.

Requires Go 1.25+ and a C compiler for SQLite.

## Run

```sh
go run .
```

Open http://127.0.0.1:8080. To choose a database or port:

```sh
go run . -db hospital.db -addr 127.0.0.1:9090
```

## Test

```sh
go test ./...
```
