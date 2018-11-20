package orderbook

type Order struct {
	ID			int			`json:id`
	Price		float64		`json:price`
	Size		float64		`json:size`
	Side		string		`json:side`
	MarketID	int			`json:market_id`
}