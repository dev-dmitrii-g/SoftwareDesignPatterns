# Biomedical signals

Variant 2. Small Go CLI using only the standard library. Go replaces the assignment's Java/JUnit stack.

```sh
go run .
go run . -workers 4 -samples 12000 -report json -target patient-1
go run . -bench -samples 1200000
go test -race ./...
```

Reports: `quality`, `alarms`, `filter`, `json`, `latency`. Targets: `hospital`, `patient-1`, `patient-2`. The preview prints five samples and unsubscribes itself. Progress goes to stderr; the report goes to stdout. Ctrl+C stops new jobs and drains accepted work. Each run starts fresh.

`main.go` wires everything. `observer.go` handles typed subscriptions. Sources push immutable sample values so listeners receive the same measurement. Callbacks use a copied listener list outside the registry lock; already captured callbacks can finish after unsubscribe.

The five sample listeners are Store, AdaptiveFilter, Inspector, Metrics and Preview. The filter publishes moving averages, and the inspector publishes patient reports. Summary receives both event types. Samples are synthetic; quality below 60 triggers a visitor on only that patient's snapshot.

`tree.go` contains Hospital → Ward → Patient → ECG/PPG/EEG/Device/Temperature. Visitors own traversal, allowing a report to inspect one subtree. Generic `Visitor[R]` supplies typed results; `Accept` supplies double dispatch. Elements hold data. JSON nesting is handwritten; only string escaping uses `encoding/json`.

The bounded worker pool defaults to three goroutines. The filter changes its window from 4 to 16 when queue depth reaches 8 or latency reaches 10 ms, reducing output frequency. Shutdown joins workers and flushes pending averages. The 31 tests include exact 4 × 1,000 delivery, cancellation draining, self-subscription, and hand-calculated visitor results.

For live coding: add a listener in `listeners.go` and register it in `NewApp`; copy `latency.go` to add a visitor, then add its report case. [Simple class diagram](uml.md).

One macOS arm64 run with Go 1.27.3, 1,200,000 samples per worker count, preview output discarded:

| Workers | Samples/s | Mean latency, ms |
| --- | ---: | ---: |
| 1 | 3,289,596 | 0.01584 |
| 2 | 3,208,041 | 0.01518 |
| 4 | 2,570,793 | 0.01817 |

Latency includes queueing and processing through inspection. Extra workers contend on shared state. The latency visitor averages the latest channel readings; the summary averages all samples.

Trade-off commits, with paths relative to this folder:

- `f757c55`: added `latency.go`, `latency_test.go`; changed `main.go`, `main_test.go`. A new operation leaves elements and the dispatcher unchanged.
- `6b64435`: changed `tree.go`, `visitors.go`, `latency.go`, `source.go`, `store.go`, `main.go`; added `temperature_test.go`. A new element requires handling in every visitor plus source wiring. Visitor works best when element types stay stable.
