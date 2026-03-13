# tossinvest-cli Funding And FX Guidance Implementation Plan

Date: 2026-03-13
Status: Drafted from approved design
Scope: prepare failure body capture, branch classification, and place guidance

## Objective

`order place`가 broker `prepare` 단계에서 `funding` 또는 `FX consent` 때문에 막히면, CLI가 원인을 분류하고 단계별 행동과 재시도 명령을 제공하게 만든다.

## Phase 1. Prepare Failure Capture

### Work

- `StatusError`가 non-2xx response body를 보존하도록 확장한다
- `place prepare` 실패 경로에서 raw body를 읽어 분류기로 전달한다
- generic unclassified prepare rejection typed error를 추가한다

### Deliverables

- updated `internal/client/errors.go`
- updated `internal/client/trading.go`
- rejection-body preservation tests

## Phase 2. Branch Classification

### Work

- `funding_required`와 `fx_consent_required` 분류기 추가
- `place prepare` 실패 시 typed branch error 반환
- ambiguous or weak signals는 `unclassified_prepare_failure`로 유지

### Deliverables

- new typed errors in `internal/trading/errors.go`
- client tests for funding, fx, and generic prepare failures

## Phase 3. Operator Guidance

### Work

- `order place`에 place-specific user-facing error formatter 추가
- branch별 1, 2, 3 단계 안내와 retry command template 출력
- README와 trading docs를 현재 guidance behavior에 맞게 갱신

### Deliverables

- updated `cmd/tossctl/errors.go`
- optional command tests for guidance rendering
- updated README and error docs

## Verification

- `go test ./...`
- `make build`
- fixture-backed tests for funding/fx classification

## Risks

- `fx_consent_required`는 캡처가 얇아서 overly broad keyword 매칭을 피해야 한다
- prepare body가 비어 있으면 guidance는 generic failure로 남는다
- 실제 브로커 응답 문구가 바뀌면 분류기를 추가 보강해야 한다

## Expected Outcome

사용자는 `order place`가 막혔을 때 다음을 즉시 알 수 있다.

- 지금 부족한 게 잔액인지, 환전/외화 사용 동의인지
- 앱/웹에서 무엇을 먼저 해야 하는지
- 완료 후 어떤 `preview`/`place` 명령으로 다시 시도해야 하는지
