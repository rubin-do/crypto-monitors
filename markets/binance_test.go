package markets

import "testing"

func TestGetBinancePrice(t *testing.T) {
	prices := make(chan float64)

	go MonitorBinancePrice(prices)

	price, ok := <-prices

	if !ok || price < 0. || price > 1000. {
		t.Errorf("Failed getting binance price, got %g!", price)
	}

	t.Logf("Got price: %g", price)
}
