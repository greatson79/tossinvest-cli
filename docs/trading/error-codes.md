# Trading Error Codes

This file will map Toss Securities trading rejection and failure signals to normalized CLI errors.

Current status:

- no broker rejection bodies captured yet
- only preview-stage signals observed

## Current Known Failure Classes

### Session Failures

Already normalized in the CLI:

- no active session
- stored session rejected

These remain valid for trading mode as well.

### Preview and Entry Preconditions

Observed endpoints that likely gate preview and submit:

- `GET /api/v2/trading/order/{stockCode}/prerequisite`
- `GET /api/v3/trading/order/{stockCode}/trading-status`
- `GET /api/v1/trading/orders/calculate/{stockCode}/orderable-quantity/sell`
- `GET /api/v2/trading/orders/calculate/{stockCode}/cost-basis-elements`

Expected future normalized failures from these classes:

- market not tradable
- account not eligible
- insufficient sellable quantity
- unsupported order type for current market
- product-risk acknowledgement required

### Buying-Power and Quantity Signals

Observed UI signals:

- buy-side page displayed `구매가능 금액 14원`
- sell-side quantity field displayed `최대 114주 가능`

Expected future normalized failures:

- insufficient buying power
- quantity exceeds sellable shares
- fractional mode mismatch

### Product Risk Gates

Observed on TSLL:

- clicking `구매하기` opened a leveraged or inverse ETP risk notice dialog
- no submit mutation was observed before the dialog

Expected normalized class:

- `product_ack_required`

CLI implication:

- preview and preflight must be able to surface product-risk acknowledgement requirements separately from broker rejection

### Ambiguous Submit

Not yet captured, but must map to a separate error family:

- request committed but response missing
- timeout after submit
- connection drop during mutation

CLI rule:

- never auto-retry on ambiguous submit
- always reconcile through order-history fetch

## To Capture Next

- rejection body for an invalid live order
- any structured broker error codes
- submit-time validation vs preflight validation differences
