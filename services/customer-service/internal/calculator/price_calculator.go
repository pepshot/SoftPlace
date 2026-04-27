package calculator

type PriceItem struct {
	Price float64
	Count int
}

func CalculateTotalPrice(items []PriceItem) float64 {
	var total float64

	for _, item := range items {
		total += item.Price * float64(item.Count)
	}

	return total
}
