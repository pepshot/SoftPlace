package helpers

import "github.com/google/uuid"

type StockItem struct {
	ID    uuid.UUID
	Count int
}

func CalculateStockDiff(oldItems []StockItem, newItems []StockItem) (toReserve []StockItem, toRelease []StockItem) {
	oldMap := make(map[uuid.UUID]int)
	newMap := make(map[uuid.UUID]int)

	for _, item := range oldItems {
		oldMap[item.ID] += item.Count
	}

	for _, item := range newItems {
		newMap[item.ID] += item.Count
	}

	for id, newCount := range newMap {
		oldCount := oldMap[id]

		if newCount > oldCount {
			toReserve = append(toReserve, StockItem{
				ID:    id,
				Count: newCount - oldCount,
			})
		}
	}

	for id, oldCount := range oldMap {
		newCount := newMap[id]

		if oldCount > newCount {
			toRelease = append(toRelease, StockItem{
				ID:    id,
				Count: oldCount - newCount,
			})
		}
	}

	return toReserve, toRelease
}
