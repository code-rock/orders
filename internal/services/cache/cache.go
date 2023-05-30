package cache

import (
	streaming "basket/internal/services/nuts-streaming"
)

type Cache struct {
	orders map[string]streaming.SOrder
}

func (c Cache) Get(id string) (order streaming.SOrder) {
	return c.orders[id]
}

func (c Cache) Set(new_order streaming.SOrder) {
	c.orders[new_order.OrderUID] = new_order
}
