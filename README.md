# toss-investment-cli

Unofficial CLI for Toss Securities web data.

## Status

This repository is in bootstrap stage. The current codebase provides:

- a Go-first CLI skeleton
- browser-assisted login and reusable session storage
- read-only account, portfolio, orders, watchlist, and quote commands
- reverse-engineering docs and sanitized fixtures
- an early trading command surface stub for future work

Live trading is not implemented yet. The current working functionality is still read-only.

## Architecture

- `Go`: main CLI, domain model, read-only client, output rendering, session lifecycle
- `Python`: future browser login helper and reverse-engineering utilities
- `Rust`: optional later addition for isolated performance-sensitive workers if real need appears

Tracked design references live in:

- [`docs/reverse-engineering/`](docs/reverse-engineering/)
- `docs/trading/` once trading discovery begins

## Current Command Surface

```bash
tossctl auth login
tossctl auth status
tossctl auth logout

tossctl account list
tossctl account summary
tossctl portfolio positions
tossctl portfolio allocation
tossctl orders list
tossctl watchlist list
tossctl quote get <symbol>
tossctl order preview
tossctl order place
tossctl order cancel
tossctl order amend
tossctl order permissions status
tossctl export positions --format csv
tossctl export orders --format json
```

`auth login`, `auth import-playwright-state`, `auth status`, `auth logout`, `quote get <symbol>`, `account list`, `account summary`, `orders list`, `portfolio positions`, `portfolio allocation`, and `watchlist list` work today.

`order` commands are present as stubs only. They define the future trading surface but do not execute mutations yet.

`auth status` performs a live validation check when a stored session exists. Authenticated commands return a re-login prompt when the stored session is missing or rejected.

## Local Paths

By default, the CLI uses OS-native paths:

- config dir: `$(os.UserConfigDir)/tossctl`
- cache dir: `$(os.UserCacheDir)/tossctl`
- session file: `<config dir>/session.json`

During development you can override paths with:

- `--config-dir`
- `--session-file`

## Development

```bash
make tidy
make fmt
make build
make test
cd auth-helper && python3 -m pip install -e . && python3 -m playwright install chromium
./bin/tossctl --help
./bin/tossctl auth login
./bin/tossctl auth status
./bin/tossctl quote get A005930
./bin/tossctl account list
./bin/tossctl account summary
./bin/tossctl portfolio positions
./bin/tossctl watchlist list --output json
```

## Safety Boundary

Today the project is read-only in practice. Trading design is underway, but any future mutation flow will stay gated behind explicit danger approvals and separate internal modules.

## Warning

This project is unofficial and not affiliated with Toss Securities. Internal web APIs can change or break without notice. Use it only if you understand those risks.
