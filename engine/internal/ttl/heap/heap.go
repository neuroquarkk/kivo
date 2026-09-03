package heap

import "container/heap"

type Heap struct {
	pq priorityQueue
}

func New() *Heap {
	return &Heap{
		pq: make(priorityQueue, 0),
	}
}

func (h *Heap) Push(key string, expiresAt int64) *Item {
	item := &Item{
		Key:       key,
		ExpiresAt: expiresAt,
	}
	heap.Push(&h.pq, item)
	return item
}

func (h *Heap) Update(item *Item, expiresAt int64) {
	item.ExpiresAt = expiresAt
	heap.Fix(&h.pq, item.Index)
}

func (h *Heap) Remove(item *Item) {
	heap.Remove(&h.pq, item.Index)
}

func (h *Heap) Pop() *Item {
	if len(h.pq) == 0 {
		return nil
	}
	return heap.Pop(&h.pq).(*Item)
}

func (h *Heap) Peek() *Item {
	if len(h.pq) == 0 {
		return nil
	}
	return h.pq[0]
}
