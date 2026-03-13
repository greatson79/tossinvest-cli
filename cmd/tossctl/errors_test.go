package main

import (
	"strings"
	"testing"

	"github.com/junghoonkye/tossinvest-cli/internal/config"
	"github.com/junghoonkye/tossinvest-cli/internal/trading"
)

func TestUserFacingPlaceErrorFormatsFundingGuidance(t *testing.T) {
	t.Parallel()

	err := userFacingPlaceError(rootPathsForTest(), &trading.BranchRequiredError{
		Branch:        trading.BranchFundingRequired,
		StatusCode:    422,
		BrokerMessage: "계좌 잔액이 부족해요 / 구매를 위해 21,511원을 채울게요.",
	}, &placeFlags{
		symbol:       "TSLL",
		market:       "us",
		side:         "buy",
		orderType:    "limit",
		quantity:     1,
		price:        500,
		currencyMode: "KRW",
	})
	if err == nil {
		t.Fatal("expected formatted error")
	}

	message := err.Error()
	if !strings.Contains(message, "잔액 또는 주문가능금액이 부족") {
		t.Fatalf("expected funding guidance, got %q", message)
	}
	if !strings.Contains(message, "tossctl order preview --symbol TSLL --market us --side buy --type limit --qty 1 --price 500 --currency-mode KRW") {
		t.Fatalf("expected preview retry command, got %q", message)
	}
	if !strings.Contains(message, "--confirm <new-confirm-token>") {
		t.Fatalf("expected execute retry template, got %q", message)
	}
}

func TestUserFacingPlaceErrorFormatsFXGuidance(t *testing.T) {
	t.Parallel()

	err := userFacingPlaceError(rootPathsForTest(), &trading.BranchRequiredError{
		Branch:        trading.BranchFXConsentRequired,
		StatusCode:    500,
		BrokerMessage: "환전 후 주문하려면 외화 사용 동의가 필요해요.",
	}, &placeFlags{
		symbol:       "TSLL",
		market:       "us",
		side:         "buy",
		orderType:    "limit",
		quantity:     1,
		price:        500,
		currencyMode: "KRW",
	})
	if err == nil {
		t.Fatal("expected formatted error")
	}

	message := err.Error()
	if !strings.Contains(message, "환전 또는 외화 사용 동의가 필요") {
		t.Fatalf("expected fx guidance, got %q", message)
	}
	if !strings.Contains(message, "Toss 앱 또는 웹에서 해당 미국주식 주문의 환전 또는 외화 사용 동의 화면으로 이동") {
		t.Fatalf("expected fx steps, got %q", message)
	}
}

func rootPathsForTest() config.Paths {
	return config.Paths{}
}
