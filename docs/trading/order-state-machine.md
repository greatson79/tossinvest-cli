# Trading Order State Machine

Initial observations captured from the TSLL order page on 2026-03-11.

## Current Working Model

The web order page appears to build trading state in layers.

### 1. Enter Order Page

Transition:

- authenticated account page
- stock order page `/stocks/US20220809012/order`

Observed network:

- trading status
- prerequisite
- pending/completed order history
- sellable quantity
- cost-basis elements
- buy-side `order-data` calculations

Interpretation:

- the page assembles enough state to preview both buy and sell before the user submits anything

### 2. Buy Preview Ready

Observed UI:

- `구매`
- `지정가`
- default price prefilled
- quantity empty

Observed side effects:

- two `order-data` POST calculations during load

Interpretation:

- the page primes a default buy preview without waiting for explicit user entry

### 3. Quantity Entered

Observed UI after entering quantity:

- total order amount recalculated
- buy-after estimate updated locally

Observed result:

- UI changed immediately
- no additional preview request was isolated from quantity change alone

Interpretation:

- either quantity recalculation is local after page bootstrapping
- or the relevant network trigger is tied to a different event boundary than simple fill/blur

### 4. Sell Preview Ready

Observed transition:

- switching from `구매` to `판매`
- setting quantity to `1`

Observed UI:

- placeholder changed to `최대 114주 가능`
- summary fields changed to:
  - `현재 수익`
  - `예상 수익률`
  - `예상 손익`
  - `총 금액`

Observed network:

- prerequisite re-fetch
- cost-basis-elements fetch

Interpretation:

- sell mode depends on cost basis and position data more explicitly than buy mode

### 5. Pre-Submit Product Risk Gate

Observed transition:

- buy mode
- quantity set to `1`
- `구매하기` clicked

Observed result:

- no submit mutation captured
- a leveraged or inverse ETP disclosure dialog blocked further progress

Interpretation:

- some products insert an acknowledgement gate between preview and submit
- the submit state machine likely includes:
  - preview-ready
  - product-risk-ack-required
  - confirmation
  - submit

## Provisional State Graph

```text
entered-page
  -> buy-preview-bootstrapped
  -> buy-quantity-updated
  -> sell-preview-bootstrapped
  -> submit-confirmation? (not yet captured)
  -> submitted? (not yet captured)
```

## Not Yet Captured

- confirmation modal state
- submit in-flight state after product-risk acknowledgement
- successful pending order state transition
- rejection transition
- cancel transition
- amend transition
- ambiguous timeout transition
