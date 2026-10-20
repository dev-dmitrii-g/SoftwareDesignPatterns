# Hospital Admission Strategies

Go console demo with two independent Strategy families. All admission outcomes and coverage records are simulated.

## Run

```sh
go run ./cmd
go test ./...
```

## Classes

| Type | Responsibility |
| --- | --- |
| `admission.Admission` | Holds the details needed to begin an admission. |
| `admission.Result` | Holds the outcome of an admission flow. |
| `flow.Strategy` | Defines the admission flow operation. |
| `flow.Automatic` | Completes a simulated admission immediately. |
| `flow.Manual` | Requests staff approval. |
| `flow.Scheduled` | Queues an admission for a supplied time. |
| `flow.Context` | Runs its current admission strategy. |
| `flow.Selector` | Selects a flow from a name, urgency, appointment time. |
| `flow.Processor` | Keeps the original API used by the refactoring tests. |
| `coverage.Request` | Holds the patient, service, evaluation time. |
| `coverage.Decision` | Holds the selected program's approval result. |
| `coverage.Record` | Holds simulated eligibility data. |
| `coverage.Source` | Defines how to look up a coverage record. |
| `coverage.MemorySource` | Looks up coverage records in memory. |
| `coverage.Strategy` | Defines the coverage check operation. |
| `coverage.GOBMP` | Checks public service eligibility. |
| `coverage.OSMS` | Checks active coverage plus referral rules. |
| `coverage.DMS` | Checks policy expiry plus covered services. |
| `coverage.Context` | Runs its current coverage strategy. |
| `coverage.Selector` | Selects a program by name or eligibility priority. |

## Why Strategy

The program has two operations with several possible algorithms. Admission can begin automatically, wait for staff approval, or be scheduled. Coverage can use GOBMP, OSMS, or DMS rules. Each algorithm implements one small interface. Its context receives an implementation in the constructor and calls it through that interface. `SetStrategy` replaces the implementation while the program runs. The context needs no branch for a particular algorithm. Selectors contain the manual and automatic choice rules. A new algorithm can be implemented without changing either context.

This is Strategy rather than State because the caller chooses the algorithm for a request. The returned admission status is an outcome, not an object controlling a patient's lifecycle. Replacing Manual with Automatic changes the algorithm used by the context; it does not move one admission through a series of states.

This is Strategy rather than Template Method because there is no parent class defining a fixed sequence of steps. Every concrete strategy implements its entire operation. The program composes a context with an interchangeable implementation at runtime. Go interfaces supply that common contract without inheritance.

## Refactoring commits

- [Before: admission modes in one switch](https://github.com/dev-dmitrii-g/SoftwareDesignPatterns/commit/28ec051)
- [After: admission flow strategies](https://github.com/dev-dmitrii-g/SoftwareDesignPatterns/commit/7c380ba)

The same processor tests pass on both commits. Go table-driven tests cover every strategy and selector rule in the final project.
