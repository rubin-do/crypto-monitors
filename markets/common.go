package markets

type Order struct {
	ItemId         string
	Market         string
	SellerName     string
	BuyPrice       float64
	SellPrice      float64
	Quantity       string
	MinAmount      string
	MaxAmount      string
	PaymentMethods string
	Url            string
}
