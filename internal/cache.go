package internal

import "sync"

var cache = struct {
	sync.RWMutex
	data map[string]Order
}{data: make(map[string]Order)}

func SaveToCache(order Order) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[order.OrderUID] = order
}

func GetFromCache(orderUID string) (Order, bool) {
	cache.RLock()
	defer cache.RUnlock()
	order, ok := cache.data[orderUID]
	return order, ok
}
