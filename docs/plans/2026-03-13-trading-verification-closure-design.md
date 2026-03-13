# tossinvest-cli Trading Verification Closure Design

Date: 2026-03-13
Status: Approved
Scope: US buy limit / KRW / non-fractional trading beta only

## Goal

현재 지원 중인 거래 베타를 확장하지 않고, live evidence 기준으로 신뢰 가능한 상태로 닫는다.

이번 단계의 목적은 새 주문 유형을 여는 것이 아니라, 이미 구현된 `place`, `cancel`, `amend`, `orders completed`, `order show <id>` 흐름이 실제로 어느 범위까지 믿을 수 있는지 증거와 문서로 고정하는 것이다.

## In Scope

- `US buy limit / KRW / non-fractional` 슬라이스 재검증
- `order place`, `order cancel`, `order amend` live reverification
- mutation 이후 상태 가시성 점검
  - `accepted_pending`
  - `filled_completed`
  - `canceled`
  - `amended_pending`
  - `amended_completed`
  - `unknown`
- 실제 검증 결과를 기준으로 README 지원 범위와 운영 문서 갱신
- 검증 중 확인된 실패 분기만 에러 모델에 반영

## Out of Scope

- `sell`
- `market`
- `KR`
- `fractional`
- 브라우저-assisted challenge 우회 자동화
- 새 주문 명령 추가

## Why This Milestone

현재 저장소는 거래 hardening의 주요 코드가 이미 들어와 있다. 남은 공백은 구현 자체보다 "무엇이 실제로 검증됐는가"에 더 가깝다.

따라서 이번 단계는 기능 추가보다 verification-first 접근이 맞다. 먼저 live evidence를 만들고, 그 결과에 맞춰 문서와 사용자 메시지를 닫아야 README와 CLI가 과장되지 않는다.

## Approach Options

### 1. Verification-first Closure

먼저 현재 구현된 명령을 live 기준으로 재검증하고, 그 결과로 문서와 메시지를 정리한다.

장점:

- 현재 베타의 실제 신뢰 구간이 명확해진다
- README와 CLI 출력이 evidence-driven 상태가 된다
- 새로운 추측 기반 코드 수정이 줄어든다

단점:

- 외형상 신규 기능 추가는 거의 없다

### 2. Error-model-first

interactive auth, funding, FX consent, product acknowledgement를 먼저 세분화한다.

장점:

- 실패 UX를 빨리 개선할 수 있다

단점:

- fresh capture 없이 들어가면 실제 분기와 어긋날 수 있다

### 3. Coverage-first Pilot

`US sell`이나 `US market` 같은 새 슬라이스를 먼저 연다.

장점:

- 지원 표면이 빨리 넓어진다

단점:

- 현재 마일스톤의 목적이 흐려진다
- 검증 매트릭스만 먼저 커진다

## Recommendation

`Verification-first Closure`로 간다.

지금 필요한 것은 breadth보다 confidence다. 현재 구현이 실제로 어떤 상태로 귀결되는지 먼저 확인하고, 그 결과로 문서와 오류 표면을 닫는 편이 맞다.

## Execution Shape

### 1. Verification Phase

현재 구현된 명령을 그대로 사용해 실제 동작을 다시 검증한다.

대상 명령:

- `order place`
- `order cancel`
- `order amend`
- `orders completed`
- `order show <id>`

이 단계에서는 코드를 먼저 넓히지 않는다. 각 명령이 어떤 입력에서 어떤 최종 상태로 보이는지 evidence를 수집하는 것이 우선이다.

### 2. Closure Phase

검증 결과를 바탕으로 필요한 수정만 한다.

범위:

- README 지원 매트릭스 정리
- live verification note 정리
- 실제 확인된 `unknown`, challenge, funding, FX, product-ack 분기만 사용자 메시지에 반영
- 필요 시 소규모 parser, output, error mapping 수정

## Verification Criteria

### order place

기대 결과:

- `accepted_pending`
- `filled_completed`
- `unknown`

### order cancel

기대 결과:

- `canceled`
- 또는 실제 미해결 시 명확한 후속 안내

### order amend

기대 결과:

- `amended_pending`
- `amended_completed`
- `unknown`

### orders completed

기대 결과:

- 방금 실행한 주문이 completed history에서 확인 가능

### order show <id>

기대 결과:

- pending/completed 중 어디에 있는지 단일 진입점으로 확인 가능

## Deliverables

- live verification note 1건
- README 업데이트
- 필요 시 소규모 코드 수정

## Success Criteria

- 현재 지원 슬라이스에 대해 `place`, `cancel`, `amend` 각각의 검증 결과가 문서로 남아 있다
- 사용자는 실행 후 상태를 `orders completed` 또는 `order show <id>`로 재확인할 수 있다
- README와 CLI 메시지가 실제 검증 결과를 넘어서 주장하지 않는다
- 남은 리스크가 "무엇이 미검증인지" 형태로 명시된다
