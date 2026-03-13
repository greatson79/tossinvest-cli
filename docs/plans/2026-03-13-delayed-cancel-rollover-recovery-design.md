# tossinvest-cli Delayed Cancel Rollover Recovery Design

Date: 2026-03-13
Status: Approved
Scope: delayed completed-history rollover recovery for `order show <old-id>`

## Goal

`cancel` mutation 직후에는 completed-history row가 아직 보이지 않아서 `current_order_id`를 못 기록하는 경우가 있다. 이번 단계의 목표는 같은 로컬 `config dir`에서 실행한 그런 주문에 대해, 나중에 `order show <old-id>`가 completed history를 다시 조회해 surviving ref를 회복하도록 만드는 것이다.

이번 단계는 background sync를 추가하지 않는다. recovery는 `order show` 요청 시점에만 수행한다.

## In Scope

- local lineage cache에 unresolved cancel lookup hint 저장
- `order show`에서 lineage cache miss 이후 on-demand completed-history fallback 수행
- recovered ref를 lineage cache에 다시 기록
- README와 verification note를 delayed rollover recovery 기준으로 갱신

## Out of Scope

- background polling or daemon-style lineage sync
- `dangerous_automation` runtime handler
- `amend` interactive auth 자동화
- cross-machine lineage sharing

## Approach Options

### 1. Cache-only

mutation 시점에 잡은 `original -> current`만 저장하고, `order show`는 그 값만 사용한다.

장점:

- 구현이 가장 단순하다

단점:

- delayed completed-history rollover를 못 잡는다
- 방금 live retest에서 확인한 실제 gap이 그대로 남는다

### 2. Cache + On-demand History Fallback

lineage cache에 unresolved cancel hint를 남기고, `order show <old-id>`가 실패할 때 completed history를 다시 조회해 surviving ref를 추정한다. 단일 candidate면 cache를 갱신하고 그 ref를 반환한다.

장점:

- delayed history 반영을 흡수할 수 있다
- background task 없이 CLI 요청 한 번으로 복구된다
- same-machine lineage 모델을 그대로 유지한다

단점:

- matching 규칙을 보수적으로 잡아야 한다
- ambiguous candidate가 생기면 자동 복구를 포기해야 한다

### 3. Background Sync

mutation 이후 별도 폴링으로 lineage file을 나중에 채운다.

장점:

- `order show`는 단순해진다

단점:

- CLI 도구에 비해 운영 복잡도가 크다
- lifecycle 관리가 불필요하게 무거워진다

## Recommendation

`Cache + On-demand History Fallback`으로 간다.

지금 필요한 것은 delayed broker history visibility를 same-machine CLI 경험 안에서 흡수하는 것이다. background sync는 과하고, cache-only는 실제로 남은 버그를 해결하지 못한다.

## Design

### 1. Lineage Cache Stores Recovery Hints

기존 lineage entry는 `current_order_id`만 저장했다. 이제 delayed cancel recovery에 필요한 최소 hint를 함께 저장한다.

- `kind`
- `symbol`
- `market`
- `quantity`
- `price`
- `order_date`
- `updated_at`

즉시 rollover를 못 잡은 cancel도 `original_order_id` 기준으로 이 hint를 로컬 파일에 남긴다.

### 2. Order Show Recovers on Demand

`order show <old-id>` 흐름은 다음 순서로 동작한다.

1. exact lookup
2. known lineage alias lookup
3. 실패 시 lineage hint 조회
4. completed history를 다시 읽어 delayed cancel candidate 탐색
5. 단일 candidate면 `original -> current` mapping을 lineage cache에 기록
6. recovered order를 반환

후보가 여러 개면 자동 추정하지 않고 ambiguous error로 실패한다.

### 3. Matching Rules Stay Conservative

completed-history fallback은 다음 조건을 모두 만족하는 row만 candidate로 본다.

- 같은 market
- 같은 symbol
- cancel로 보이는 status
- 같은 quantity
- 같은 price
- 같은 order date
- `updated_at - lookback` 이후에 보인 row
- original ref와 다른 surviving ref

candidate가 정확히 하나일 때만 recovery를 성공으로 본다.

## Deliverables

- lineage hint persistence
- delayed cancel recovery helper in completed-history lookup
- `order show` cache refresh on successful recovery
- regression tests for resolved and ambiguous cases
- updated README and verification note

## Success Criteria

- live retest에서처럼 delayed completed-history rollover가 생겨도, 같은 `config dir`에서 `order show <old-id>`가 later surviving ref를 찾을 수 있다
- recovered ref는 다음 조회부터 local lineage cache를 통해 바로 resolve된다
- candidate가 여러 개인 경우 CLI가 추정하지 않고 명확히 실패한다
- docs가 recovery 범위와 한계를 과장하지 않는다
