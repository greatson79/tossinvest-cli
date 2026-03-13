# tossinvest-cli Funding And FX Guidance Design

Date: 2026-03-13
Status: Approved
Scope: step-by-step operator guidance for funding and FX-consent prepare failures

## Goal

`order place`가 broker `prepare` 단계에서 막힐 때, CLI가 `funding_required`와 `fx_consent_required`를 구분해서 사용자가 바로 따라 할 수 있는 단계별 행동을 안내한다.

이번 단계는 자동 진행을 구현하지 않는다. 대신 이후 `dangerous_automation`으로 이어질 수 있도록 분기 모델과 사용자 메시지를 정리한다.

## In Scope

- `place prepare` 실패 응답 body 보존
- `funding_required`와 `fx_consent_required` 분류
- `order place` 전용 단계별 사용자 안내
- 재시도용 `preview`/`place` 명령 템플릿 제공
- 미분류 `prepare` 실패를 generic typed error로 래핑

## Out of Scope

- funding 자동 충전
- FX consent 자동 수락
- `product acknowledgement`
- `interactive trade auth`
- 브라우저 자동 오픈

## Approach Options

### 1. Generic Guidance

`prepare` 실패를 한 묶음으로 보고, 잔액 또는 환전 동의를 확인하라고만 안내한다.

장점:

- 구현이 가장 빠르다

단점:

- 사용자가 지금 무엇을 먼저 해야 하는지 알기 어렵다
- 나중에 자동화 분기로 연결하기도 애매하다

### 2. Response-Body Classification With Step Guide

`prepare` 실패 응답 body를 읽어 `funding_required`와 `fx_consent_required`를 구분하고, 각 분기마다 단계별 행동과 재시도 명령을 제공한다.

장점:

- 현재 요구와 정확히 맞는다
- status code 변화에 덜 흔들린다
- 이후 automation 분기 이름으로 재사용 가능하다

단점:

- 분류기와 typed error 표면이 추가된다

### 3. Heuristic Follow-up Reads

`prepare` 실패 후 account summary나 withdrawable/funding 관련 read API를 더 읽어 간접 추정한다.

장점:

- body가 빈약해도 힌트를 늘릴 수 있다

단점:

- 추정이 섞인다
- 첫 단계치고 구현이 과하다

## Recommendation

`Response-Body Classification With Step Guide`로 간다.

첫 단계에서 필요한 것은 정확한 사람 안내다. 자동화나 추정 강화보다, broker가 지금 무엇 때문에 막았는지 분리해서 바로 다음 행동을 명확히 주는 것이 우선이다.

## Design

### 1. Preserve Prepare Failure Bodies

현재 `postJSONBytes`는 non-2xx response body를 버린다. 이를 바꿔서 `StatusError`가 status, endpoint, raw body를 함께 보존하게 만든다.

이렇게 해야 `place prepare` 단계에서 broker rejection의 title, body, message, action label을 분류에 사용할 수 있다.

### 2. Classify Only What We Can Defend

`place prepare` 실패 분류기는 body 안의 문구를 기준으로 다음만 구분한다.

- `funding_required`
- `fx_consent_required`
- `unclassified_prepare_failure`

`funding_required`는 `계좌 잔액이 부족해요`, `채울게요`, `모바일에서 채우기` 같은 문구를 우선 사용한다.

`fx_consent_required`는 `환전`, `외화 사용`, `동의`, `consent` 같이 FX approval을 직접 가리키는 문구가 충분히 있을 때만 반환한다.

분류가 모호하면 추정하지 않고 `unclassified_prepare_failure`로 남긴다.

### 3. Place-Specific Operator Guidance

`order place`는 일반 trading error formatter 대신 place 전용 formatter를 거쳐, 다음 구조의 안내를 보여준다.

- 원인 한 줄
- 지금 할 일 1, 2, 3
- 다시 실행할 `preview` 명령
- 새 confirm token으로 다시 실행할 `place` 명령 템플릿

이 formatter는 현재 place flags를 알고 있으므로, 실제 입력값이 들어간 명령 예시를 그대로 보여줄 수 있다.

## Deliverables

- prepare rejection body preservation
- funding/fx typed errors
- place-specific step-by-step guidance formatter
- client tests for classification
- command tests for user guidance
- updated README and trading docs

## Success Criteria

- 사용자가 CLI 출력만 보고 `funding`과 `FX consent`를 구분할 수 있다
- 출력에 단계별 행동과 재시도 명령이 포함된다
- 분류가 불충분할 때는 과감히 generic prepare failure로 남는다
- 이후 automation 구현에서 같은 branch name을 재사용할 수 있다
