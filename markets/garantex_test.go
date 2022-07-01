package markets

import "testing"

func TestGetGarantexOrder(t *testing.T) {
	orders := make(chan Order)

	go MonitorGarantexPrice(orders)

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
