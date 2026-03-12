package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/junghoonkye/tossinvest-cli/internal/domain"
)

type completedOrdersEnvelope struct {
	Result struct {
		Body []json.RawMessage `json:"body"`
	} `json:"result"`
}

func (c *Client) ListCompletedOrders(ctx context.Context, market string) ([]domain.Order, error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return c.ListCompletedOrdersRange(ctx, market, from, now, 50, 1)
}

func (c *Client) ListCompletedOrdersRange(ctx context.Context, market string, from, to time.Time, size, number int) ([]domain.Order, error) {
	if err := c.requireSession(); err != nil {
		return nil, err
	}
	if err := c.ensureTradingMetadata(ctx); err != nil {
		return nil, err
	}

	markets, err := normalizeHistoryMarkets(market)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.Order, 0)
	for _, entry := range markets {
		var envelope completedOrdersEnvelope
		endpoint := fmt.Sprintf(
			"%s/api/v2/trading/my-orders/markets/%s/by-date/completed?range.from=%s&range.to=%s&size=%d&number=%d",
			c.certBaseURL,
			entry,
			from.Format("2006-01-02"),
			to.Format("2006-01-02"),
			size,
			number,
		)
		if err := c.getJSON(ctx, endpoint, &envelope); err != nil {
			return nil, err
		}

		for _, item := range envelope.Result.Body {
			orders = append(orders, parseCompletedOrder(item, entry))
		}
	}

	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].SubmittedAt == nil && orders[j].SubmittedAt == nil {
			return orders[i].ID > orders[j].ID
		}
		if orders[i].SubmittedAt == nil {
			return false
		}
		if orders[j].SubmittedAt == nil {
			return true
		}
		return orders[i].SubmittedAt.After(*orders[j].SubmittedAt)
	})

	return orders, nil
}

func (c *Client) FindOrder(ctx context.Context, orderID string, market string) (domain.Order, error) {
	pendingOrders, err := c.ListPendingOrders(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	for _, order := range pendingOrders {
		if order.ID == orderID || orderMatchesID(order.Raw, orderID) {
			return order, nil
		}
	}

	completedOrders, err := c.ListCompletedOrders(ctx, market)
	if err != nil {
		return domain.Order{}, err
	}
	for _, order := range completedOrders {
		if order.ID == orderID || orderMatchesID(order.Raw, orderID) {
			return order, nil
		}
	}

	return domain.Order{}, fmt.Errorf("order %s was not found in pending or current-month completed history", orderID)
}

func parseCompletedOrder(raw json.RawMessage, market string) domain.Order {
	order := domain.Order{Raw: raw}

	var payload struct {
		OrderedAt      string `json:"orderedAt"`
		LastExecutedAt string `json:"lastExecutedAt"`
		OrderNo        any    `json:"orderNo"`
		OrderID        string `json:"orderId"`

		StockCode string `json:"stockCode"`
		StockName string `json:"stockName"`
		Symbol    string `json:"symbol"`
		TradeType string `json:"tradeType"`
		Status    string `json:"status"`

		OrderQuantity    float64 `json:"orderQuantity"`
		ExecutedQuantity float64 `json:"executedQuantity"`
		UserOrderDate    string  `json:"userOrderDate"`

		OrderPrice struct {
			KRW float64 `json:"krw"`
		} `json:"orderPrice"`
		AverageExecutionPrice struct {
			KRW float64 `json:"krw"`
		} `json:"averageExecutionPrice"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return order
	}

	order.ID = referenceOrderIdentifier(payload.UserOrderDate, payload.OrderNo, payload.OrderID)
	order.Symbol = firstNonEmpty(payload.Symbol, payload.StockCode)
	order.Name = payload.StockName
	order.Market = strings.ToLower(market)
	order.Side = payload.TradeType
	order.Status = payload.Status
	order.Quantity = payload.OrderQuantity
	order.FilledQuantity = payload.ExecutedQuantity
	order.Price = payload.OrderPrice.KRW
	order.AverageExecutionPrice = payload.AverageExecutionPrice.KRW
	order.OrderDate = payload.UserOrderDate
	order.SubmittedAt = parseOrderTime(payload.LastExecutedAt, payload.OrderedAt)
	return order
}

func normalizeHistoryMarkets(value string) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "all":
		return []string{"us", "kr"}, nil
	case "us", "kr":
		return []string{strings.ToLower(strings.TrimSpace(value))}, nil
	default:
		return nil, fmt.Errorf("unsupported market %q; use us, kr, or all", value)
	}
}
