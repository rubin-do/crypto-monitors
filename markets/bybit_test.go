package markets

import (
	"testing"
)

func TestGetByBitOrder(t *testing.T) {
	orders := make(chan Order)

	go MonitorByBitPrice(orders)

	order, ok := <-orders
	price := order.BuyPrice

	if !ok || price <= 0. || price > 1000. {
		t.Errorf("Failed getting binance price, got %g!", price)
	}

	t.Logf("Username: %s\nBuyPrice: %g\nSellPrice: %g\nQuantity: %s\nMin: %s\nMax: %s\nUrl: %s\n",
		order.SellerName,
		order.BuyPrice,
		order.SellPrice,
		order.Quantity,
		order.MinAmount,
		order.MaxAmount,
		order.Url,
	)
}
