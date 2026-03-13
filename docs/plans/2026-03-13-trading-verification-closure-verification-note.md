# tossinvest-cli Trading Verification Closure Verification Note

Date: 2026-03-13
Status: In progress
Scope: US buy limit / KRW / non-fractional only

## Completed Checks

### Automated

- `go test ./...`
  - result: pass
- `make build`
  - result: pass

### Safe CLI Readiness

- `tossctl version`
  - result: pass
- `tossctl doctor`
  - result: pass
  - session file exists
  - trading permission file exists but temporary permission is expired
  - config file does not exist yet, so trading actions default to disabled
- `tossctl auth doctor`
  - result: pass
  - auth helper importable
  - playwright installed
  - chromium installed
  - stored session valid
- `tossctl auth status`
  - result: active session
  - provider: `playwright-storage-state`
  - live check: valid

### Read-only Order Visibility

- `tossctl orders list --output json`
  - result: pass
  - observed state: no pending orders at the time of check
- `tossctl orders completed --market us --output json`
  - result: pass
  - observed state: completed-history lookup works for current month US orders
  - observed statuses in history: `체결완료`, `취소`, `실패`
- `tossctl order show 2026-03-11/25 --market us --output json`
  - result: pass
  - observed state: canceled order lookup works through the single-order surface
- `tossctl order show 2026-03-11/1 --market us --output json`
  - result: pass
  - observed state: completed order lookup works through the single-order surface
- `tossctl order preview --symbol TSLL --market us --side buy --type limit --qty 1 --price 500 --currency-mode KRW --output json`
  - result: pass
  - observed state: preview emits canonical intent and confirm token
  - observed state: `live_ready=true`, `mutation_ready=false` while config remains disabled

## Current Blockers for Live Mutation Verification

- `config.json` does not exist yet, so:
  - `place=false`
  - `cancel=false`
  - `amend=false`
  - `allow_dangerous_execute=false`
- temporary trading permission is expired
- any live verification of `order place`, `order amend`, or `order cancel` would affect the real account and should not be run implicitly

## Still Pending

- live `order place` verification
- post-place `orders completed` and `order show <id>` verification against the new order
- live `order amend` verification
- live `order cancel` verification
- evidence-driven README updates after live verification
- any code or message corrections found during live verification

## Next Operator Steps

1. Decide whether to proceed with live account verification.
2. If yes, initialize `config.json` and explicitly enable only the required actions.
3. Refresh trading permission with `tossctl order permissions grant --ttl ...`.
4. Run the live verification sequence from the implementation plan.
5. Record the exact returned status and order ids for each mutation.
