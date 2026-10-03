# Hospital Admission

A simple Go and SQLite application that uses a facade to coordinate insurance verification, equipment reservation, and pharmacy dispatch. If a later step fails, the facade reverses completed steps. The website provides Admissions and Stock pages.

Requires Go 1.25+ and a C compiler for SQLite.

## Run the website

```sh
go run . -db hospital.db serve
```

Open http://127.0.0.1:8080.

## CLI

```sh
go run . -db hospital.db admit A001 990101123456 OSMS wheelchair paracetamol
go run . -db hospital.db inspect A001
go run . -db hospital.db stock equipment wheelchair 0
go run . -db hospital.db stock pharmacy paracetamol 0
```

## Test

```sh
go test ./...
```
