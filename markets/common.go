package markets

type Order struct {
	Market     string
	SellerName string
	BuyPrice   float64
	SellPrice  float64
	Quantity   string
	MinAmount  string
	MaxAmount  string
	Url        string
}
