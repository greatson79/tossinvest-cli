# Trading Discovery

This directory is the tracked home for trading-specific reverse-engineering notes.

It is intentionally separate from the local-only planning documents under `docs/plans/`.

Expected contents:

- `rpc-catalog.md`
- `order-state-machine.md`
- `error-codes.md`

Rules:

- do not commit raw storage-state files
- do not commit raw order captures
- do not commit secrets, tokens, or account numbers
- sanitize trading responses before adding them here
