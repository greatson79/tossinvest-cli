# tossinvest-cli Delayed Cancel Rollover Recovery Implementation Plan

Date: 2026-03-13
Status: Drafted from approved design
Scope: local recovery hints, on-demand completed-history lookup, and docs

## Objective

즉시 reconciliation에서 `current_order_id`를 못 잡은 cancel에 대해서도, 같은 머신의 `order show <old-id>`가 delayed completed-history row를 찾아 surviving ref를 복구하게 만든다.

## Phase 1. Recovery Hint Persistence

### Work

- lineage entry에 delayed recovery용 hint 필드를 추가한다
- mutation 결과를 lineage cache에 기록할 때 unresolved cancel도 저장한다
- direct lookup API를 추가해서 `order show`가 current alias와 recovery hint를 모두 읽을 수 있게 한다

### Deliverables

- updated `internal/orderlineage/service.go`
- lineage service tests for unresolved hint persistence

## Phase 2. On-demand Completed-History Recovery

### Work

- completed-history에서 delayed cancel successor를 찾는 helper를 추가한다
- candidate filtering을 symbol, market, quantity, price, order date, updated-at window, cancel status로 제한한다
- ambiguous candidate는 자동 복구하지 않고 explicit error를 반환한다

### Deliverables

- updated `internal/client/completed_orders.go`
- client tests for successful and ambiguous delayed recovery

## Phase 3. Order Show Integration and Docs

### Work

- `order show`가 normal lookup 실패 후 lineage hint recovery를 시도하게 한다
- recovery 성공 시 lineage cache를 refresh한다
- README와 verification note에 delayed cancel recovery behavior와 한계를 반영한다

### Deliverables

- updated `cmd/tossctl/order.go`
- updated README and verification note

## Verification

- `go test ./...`
- `make build`
- `./bin/tossctl order show <old-id>` regression path in tests

## Risks

- 같은 날짜에 동일 symbol, price, quantity를 가진 canceled row가 여러 개면 recovery는 ambiguous로 실패할 수 있다
- recovery는 같은 로컬 `config dir`에 남은 lineage hint가 있을 때만 가능하다
- `amend` delayed rollover는 이번 단계 범위가 아니다

## Expected Outcome

현재 trading beta 범위에서 delayed cancel rollover가 발생해도 사용자는 다음을 할 수 있다.

- 나중에 `order show <old-id>`로 surviving ref를 다시 찾는다
- 한 번 복구된 ref는 lineage cache를 통해 바로 resolve된다
- ambiguous case는 자동 추정하지 않고 수동 확인 경로를 안내받는다
