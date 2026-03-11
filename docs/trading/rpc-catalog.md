# Trading RPC Catalog

Verified from the TSLL order page on 2026-03-11.

Current capture scope:

- `preview-only`
- authenticated web session
- stock page: `/stocks/US20220809012/order`
- market: `us`

## Status Legend

- `observed`: seen in authenticated browser traffic
- `captured`: request and response family confirmed during a directed scenario
- `unknown`: likely relevant but not yet isolated

## Preview and Order-Entry Endpoints

| Status | Method | Host | Path | Purpose | Notes |
| --- | --- | --- | --- | --- | --- |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v3/trading/order/US20220809012/trading-status` | product trading status | loaded on page entry before order interaction |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v2/trading/order/US20220809012/prerequisite` | order-entry prerequisites | re-fetched when switching to `판매` mode |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/calculate/US20220809012/orderable-quantity/sell?forceFetch=false` | sellable quantity lookup | loaded on entry and again during sell-mode rendering |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v2/trading/orders/calculate/US20220809012/cost-basis-elements` | average cost and sell-side basis data | paired with sell-mode preview widgets |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/calculate/US20220809012/average-price?forceFetch=false` | average price lookup | used on order page initialization |
| `captured` | `POST` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/calculate/order-data` | buy-side preview calculation | two POSTs observed on page load with different `orderPrice` values |
| `observed` | `GET` | `wts-api.tossinvest.com` | `/api/v1/trading/settings/toggle/find?categoryName=TRADE_WITHOUT_CONFIRM` | fetch trade-without-confirm preference | may become relevant for permission UX mapping |
| `observed` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v2/trading/settings/investor-exchange-choice-type` | exchange-choice setting | likely affects market routing |

## Order History and Sidecar Endpoints

| Status | Method | Host | Path | Purpose | Notes |
| --- | --- | --- | --- | --- | --- |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/histories/all/pending` | global pending-order summary | loaded with the trading page shell |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/histories/PENDING?stockCode=US20220809012&number=1&size=100&marketDivision=us` | symbol-specific pending orders | shown in the order page side panel |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/orders/histories/COMPLETED?stockCode=US20220809012&number=1&size=30&marketDivision=us` | symbol-specific completed orders | loaded on entry |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v4/trading/auto-trading?productCode=US20220809012&size=20&number=1` | auto-trading metadata | likely informational, not needed for first mutation path |
| `captured` | `GET` | `wts-cert-api.tossinvest.com` | `/api/v1/trading/analysis/productCode/US20220809012` | trading analysis panel data | informational only so far |

## Captured `order-data` Request Bodies

Two buy-side preview POST bodies were captured during page load:

```json
{
  "stockCode": "US20220809012",
  "market": "us",
  "orderPrice": 0,
  "orderVolumeRate": 1,
  "currencyMode": "KRW",
  "isFractional": false
}
```

```json
{
  "stockCode": "US20220809012",
  "market": "us",
  "orderPrice": 21134,
  "orderVolumeRate": 1,
  "currencyMode": "KRW",
  "isFractional": false
}
```

Current inference:

- `orderPrice=0` is likely a market-price or placeholder calculation path.
- `orderPrice=21134` matches the default limit price shown in the buy panel.
- `orderVolumeRate=1` appears to represent full entered quantity in the current UI model.
- `currencyMode=KRW` is explicitly sent even on a US stock order page.

## UI Observations Tied To Preview

### Buy-side initial state

- mode: `구매`
- price mode: `지정가`
- price field defaulted to `21,134`
- quantity empty on initial load
- buying power shown as `14원`
- total order amount updated locally when quantity was entered

### Sell-side preview state

- mode: `판매`
- price mode: `지정가`
- price field showed `21,119`
- quantity `1` displayed placeholder `최대 114주 가능`
- UI summary showed:
  - `현재 수익`
  - `예상 수익률`
  - `예상 손익`
  - `총 금액`

Sell-mode switching triggered at least:

- `GET /api/v2/trading/order/US20220809012/prerequisite`
- `GET /api/v2/trading/orders/calculate/US20220809012/cost-basis-elements`

## Pre-Submit Gating Observations

Attempted action:

- buy mode
- quantity `1`
- clicked `구매하기`

Observed result:

- no live submit mutation was captured
- a blocking dialog appeared for leveraged or inverse ETP risk disclosure

Dialog summary:

- title family: leveraged or inverse ETP risk notice
- action button: `확인했어요`

Current inference:

- certain products require a product-risk acknowledgement step before any actual place mutation
- this gate must be modeled separately from permission flags and preview calculation

## Unknowns

- full response body shape for `order-data`
- whether sell preview also calls `order-data` under a different event boundary
- request body for actual place, cancel, and amend mutations
- idempotency or duplicate-submit protection fields
- any hidden nonce or ticket required on live submit
