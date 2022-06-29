package markets

import "testing"

func TestGetGarantexPrice(t *testing.T) {
	prices := make(chan float64)

	go MonitorGarantexPrice(prices)

	price, ok := <-prices

	if !ok || price < 0. || price > 1000. {
		t.Errorf("Failed getting binance price, got %g!", price)
	}

	t.Logf("Got price: %g", price)
}
