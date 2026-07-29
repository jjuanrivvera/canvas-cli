# Decisions

Numbered record of ambiguous design/API assumptions and non-obvious engineering
choices for canvas-cli. Each entry captures the decision, the alternatives, and
the rationale so the reasoning survives beyond the diff.

## 1. Diagnostics runner injection seam for deterministic `canvas doctor` tests

**Context.** `canvas doctor` runs environment, configuration, connectivity,
authentication, API-access, disk-space, and permission checks. Several of these
touch live network, the host filesystem, and real credentials, so the command
tests (`TestDoctorCmd`) were flaky in local runs and were skipped entirely under
`CI=1` / `GITHUB_ACTIONS` — a crutch that traded determinism for coverage.

**Decision.** Introduce a `diagnostics.Runner` interface
(`Run(ctx) (*Report, error)`, already satisfied by `*diagnostics.Doctor`) and
route the command through an injectable package var
`newDoctorRunner(cfg, client) diagnostics.Runner` in `commands/doctor.go`
instead of calling `diagnostics.New` directly. The production default is
unchanged. Command tests override the seam with a `fakeRunner` that returns a
deterministic `*diagnostics.Report`, exercising every output format (human,
human-verbose, `--json`, and global `-o json`) plus failure and runner-error
handling — with no network, host-permission, or real-credential dependency.

**Alternatives considered.**
- *Inject each check's dependency (HTTP client, filesystem) into `Doctor`.*
  Rejected as heavier than needed: the diagnostics package tests
  (`internal/diagnostics`, ~97% covered) already drive individual checks against
  `httptest` and temp dirs deterministically. The flakiness lived at the
  *command* layer, so the seam belongs there.
- *Keep the `CI=1` skip.* Rejected — it leaves the command untested in CI and
  still flaky locally, which is exactly the bug (#28).

**Consequences.**
- `go test ./commands -run TestDoctorCmd` is deterministic and host-free by
  default (≈0.02s vs ≈6s), and no longer skips under CI.
- The real end-to-end path stays covered by a single opt-in smoke test,
  `TestDoctorCmd_Live`, gated behind `CANVAS_DOCTOR_LIVE=1` and skipped by
  default.
- Runtime behavior of `canvas doctor` for real users is unchanged; this is a
  testability refactor only.

Ref: [#28](https://github.com/jjuanrivvera/canvas-cli/issues/28)
