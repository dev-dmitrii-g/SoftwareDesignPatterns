# Hospital Admission Facade

A small Go CLI demonstrates a Facade coordinating three independent services with durable SQLite state:

1. **Insurance Verifier** checks a 12-digit demo IIN and policy number, then persists `VERIFIED`.
2. **Equipment Reserver** reserves one named item from SQLite stock.
3. **Pharmacy Dispatcher** dispatches one named medicine from SQLite stock.

This is a teaching model. The IIN check only verifies 12 digits; no live Kazakhstan insurance, prescription, or hospital system is contacted.

## Run

Requires Go 1.25+, a C compiler (for `go-sqlite3`), and internet access for the first dependency download. From this directory:

```sh
go run . -db hospital.db stock equipment wheelchair 1
go run . -db hospital.db stock pharmacy paracetamol 1
go run . -db hospital.db admit A001 990101123456 POL-1 wheelchair paracetamol
go run . -db hospital.db inspect A001
```

The same `hospital.db` file is opened by every command, so state survives closing the program. An admission ID must be unique. The CLI prints the persisted status of all three subsystems.

### Demonstrate both compensation depths

Use separate IDs. These commands work on the same database used above:

```sh
# B fails: no wheelchair remains. Only insurance is reversed.
go run . -db hospital.db admit B001 990101123456 POL-2 wheelchair paracetamol
go run . -db hospital.db inspect B001

# C fails: make equipment available but medicine unavailable.
go run . -db hospital.db stock equipment wheelchair 1
go run . -db hospital.db stock pharmacy paracetamol 0
go run . -db hospital.db admit C001 990101123456 POL-3 wheelchair paracetamol
go run . -db hospital.db inspect C001
```

`admit` exits with an error for failed operations, while printing the final persisted state. B failure ends with insurance `REVERSED` and no equipment row. C failure releases equipment first, then reverses insurance; the equipment row is `RELEASED` and its stock is restored. Both admissions end `COMPENSATED` with a stored reason.

### Restart and administration

```sh
go run . -db hospital.db stock equipment wheelchair 1
go run . -db hospital.db stock pharmacy paracetamol 1
go run . -db hospital.db start R001 990101123456 POL-4 wheelchair paracetamol
# This is a new program process; it reads the stored PENDING_B status.
go run . -db hospital.db resume R001
go run . -db hospital.db inspect R001
```

`start` intentionally stops after A so the restart can be shown live. `compensate ID` manually reverses an incomplete admission. The admin facade refuses to compensate an admission whose medicine was dispatched.

To reset for a fresh demo, delete the database file (and any `hospital.db-wal` or `hospital.db-shm` files) while the app is stopped. The next command recreates it. For example: `rm -f hospital.db hospital.db-wal hospital.db-shm`.

## Design

The [UML class diagram](uml.svg) shows the concrete types. Each subsystem contains an entity with guarded status transitions, a separate business rule, its own SQLite repository, and a service. The subsystems do not import each other. Only `IntakeFacade` knows all three.

`Admit` runs A → B → C. A successful step commits before the next begins. `admissions.status` records chain progress (`PENDING_A`, `PENDING_B`, `PENDING_C`, `COMPLETED`, `COMPENSATING`, `COMPENSATED`). `Resume` reads this status and the actual subsystem rows, so it can recover if a process stops after a step commits but before the progress marker changes. On B failure, the facade reverses A. On C failure, it releases B and then reverses A. Compensation status is persisted before undoing anything, allowing a later `resume` to retry the remaining reversals.

SQLite transactions keep each stock change together with its reservation or dispatch row. This is a small single-process example, not a production concurrent admission service.

### Principle of Least Knowledge

Without the facade, a client would need several immediate collaborators and would have to know their order and compensation rules:

```go
insuranceVerifier.Verify(id, iin, policy)
if err := equipmentReserver.Reserve(id, item); err != nil {
    insuranceVerifier.Reverse(id)
}
if err := pharmacyDispatcher.Dispatch(id, medicine); err != nil {
    equipmentReserver.Release(id)
    insuranceVerifier.Reverse(id)
}
```

With the facade, the normal client has one collaborator:

```go
err := intakeFacade.Admit(admission)
```

The admin client uses `AdminFacade.Inspect` or `AdminFacade.Compensate` without reaching into individual repositories.

### Facade versus Adapter

The problem is coordinating an admission across insurance, equipment, and pharmacy. Their interfaces are already usable individually; there is no incompatible interface to translate. `IntakeFacade` provides a higher-level operation, chooses the A → B → C order, reads durable progress, and reverses committed work after a later failure. An Adapter would mainly reshape one interface for a caller and would not by itself own this chain or its compensation policy.

## Tests and coverage

```sh
go test ./...
go test -coverpkg=./internal/... -coverprofile=coverage.out ./internal/admission
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

The checked-in [coverage report](coverage.html) and `coverage.out` were generated from these commands. Coverage is **over 80%** across the three subsystems and facades. Tests use a fresh SQLite file in `t.TempDir()` for each case and assert persisted states after reopening it. The test helper reads like a domain sentence: `aHospital(t).stock(...).admits(...)`, followed by `records(...)`.

This directory is a local Git repository with source and tests committed together. Push it to your hosting account and submit that URL through Moodle; no remote repository was provided here.
