package markets

import (
	"testing"
)

func TestGetBinanceOrder(t *testing.T) {
	orders := make(chan Order)

	go MonitorBinancePrice(orders)

	order, ok := <-orders
	price := order.Price

	if !ok || price <= 0. || price > 1000. {
		t.Errorf("Failed getting binance price, got %g!", price)
	}

	t.Logf("Username: %s\nPrice: %g\nQuantity: %s\nMin: %s\nMax: %s\nUrl: %s\n",
		order.SellerName,
		order.Price,
		order.Quantity,
		order.MinAmount,
		order.MaxAmount,
		order.Url,
	)
}
