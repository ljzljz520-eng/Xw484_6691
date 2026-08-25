# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	genealogy-story-organizer/cmd/genealogy-server	0.011s
ok  	genealogy-story-organizer/internal/application	0.023s
ok  	genealogy-story-organizer/internal/domain	0.001s
--- FAIL: Test484BusinessRegression (0.01s)
    business_regression_test.go:26: second amount=210
FAIL
FAIL	genealogy-story-organizer/internal/flow023	0.013s
ok  	genealogy-story-organizer/internal/importer	0.001s
ok  	genealogy-story-organizer/internal/integration	0.013s
ok  	genealogy-story-organizer/internal/query	0.002s
ok  	genealogy-story-organizer/internal/store	0.013s
ok  	genealogy-story-organizer/internal/web	0.006s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/genealogy-server): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/genealogy-server): exit `0`
- Frontend build (web): exit `0`
