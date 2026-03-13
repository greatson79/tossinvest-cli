# tossinvest-cli Trading Verification Closure Implementation Plan

Date: 2026-03-13
Status: Drafted from approved design
Scope: US buy limit / KRW / non-fractional only

## Objective

현재 거래 베타의 구현을 넓히지 않고, live verification evidence와 문서, 최소한의 결과 보정으로 "검증된 상태"를 닫는다.

## Phase 1. Save the Verification Baseline

### Work

- 현재 설계를 `docs/plans`에 고정
- 기존 hardening 문서와 현재 코드의 차이를 명시
- 이번 단계가 coverage expansion이 아니라 verification closure라는 점을 분명히 기록

### Deliverables

- approved design doc
- implementation plan doc

### Exit Criteria

- 팀이나 에이전트가 이번 마일스톤을 기능 확장으로 오해하지 않는다

## Phase 2. Safe Readiness Checks

### Work

- `go test ./...`
- `make build`
- `tossctl version`
- `tossctl doctor`
- `tossctl auth doctor`
- 필요 시 `tossctl auth status`

### Deliverables

- repo verification 결과
- 로컬 실행 준비 상태 요약

### Exit Criteria

- live verification 전에 코드와 환경의 기본 준비 상태가 확인된다

## Phase 3. Live Verification Runbook

### Work

- 아래 시나리오를 순서대로 실행할 runbook 작성
  - `order preview`
  - `order place`
  - `orders completed`
  - `order show <id>`
  - `order amend`
  - `order cancel`
- 각 단계에서 기록할 필드 정의
  - 입력값
  - confirm token
  - returned status
  - returned order id
  - follow-up lookup result
  - mismatch or uncertainty

### Deliverables

- live verification note template 또는 runbook

### Exit Criteria

- 실제 계정 검증 시 어떤 evidence를 남겨야 하는지 모호하지 않다

## Phase 4. Closure Updates

### Work

- live verification 결과를 근거로 README 지원 범위 재정리
- 실제 확인된 실패 분기만 에러 모델과 사용자 메시지에 반영
- 필요 시 output 또는 parsing의 작은 불일치 수정

### Deliverables

- updated README
- updated verification note
- minimal code/docs fixes if required

### Exit Criteria

- README와 CLI 메시지가 evidence-driven 상태가 된다

## Verification Plan

### Automated

- `go test ./...`
- `make build`

### Non-destructive CLI

- `tossctl version`
- `tossctl doctor`
- `tossctl auth doctor`
- `tossctl auth status`

### Live Account Verification

- `tossctl order preview`
- `tossctl order place`
- `tossctl orders completed`
- `tossctl order show <id>`
- `tossctl order amend`
- `tossctl order cancel`

## Risks

- live verification은 실제 계정과 브로커 상태에 의존한다
- 일부 분기(funding, FX consent, product acknowledgement)는 같은 입력에서도 재현되지 않을 수 있다
- `unknown` 결과는 코드 문제보다 broker-side visibility delay일 수 있다

## Recommended Execution Order

1. 설계와 계획 문서 저장
2. safe readiness checks 실행
3. live verification runbook 준비
4. explicit operator confirmation 후 live account verification
5. evidence 기반 closure updates 적용

## Expected Outcome

이번 단계가 끝나면 현재 거래 베타는 "구현돼 있음"이 아니라 "무엇이 live로 검증됐고 무엇이 아직 아닌지 명확한 상태"로 정리된다.
