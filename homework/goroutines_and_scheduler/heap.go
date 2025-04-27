package main

type Heap[T any, I comparable] struct {
	data          []T
	indices       map[I]int // Identifier -> index
	less          func(T, T) bool
	getIdentifier func(T) I // item -> identifier
}

func NewHeap[T any, I comparable](less func(T, T) bool, getIdentifier func(T) I) *Heap[T, I] {
	return &Heap[T, I]{
		data:          []T{},
		indices:       map[I]int{},
		less:          less,
		getIdentifier: getIdentifier,
	}
}

func (h *Heap[T, I]) Push(item T) {
	id := h.getIdentifier(item)
	if _, ok := h.indices[id]; ok {
		return
	}

	h.data = append(h.data, item)
	lastIndex := len(h.data) - 1
	h.indices[id] = lastIndex
	h.up(lastIndex)
}

func (h *Heap[T, I]) Pop() T {
	if len(h.data) == 0 {
		return *new(T)
	}

	result := h.data[0]
	last := h.data[len(h.data)-1]

	h.data = h.data[:len(h.data)-1]
	delete(h.indices, h.getIdentifier(result))

	if len(h.data) > 0 {
		h.data[0] = last
		h.indices[h.getIdentifier(last)] = 0
		h.down(0)
	}

	return result
}

func (h *Heap[T, I]) Change(identifier I, update func(T) T) {
	index, ok := h.indices[identifier]
	if !ok {
		return
	}

	item := h.data[index]
	h.data[index] = update(item)

	h.down(index)
	h.up(index)
}

func (h *Heap[T, I]) GetByIdentifier(identifier I) (T, bool) {
	index, ok := h.indices[identifier]
	if !ok {
		return *new(T), false
	}

	return h.data[index], true
}

func (h *Heap[T, I]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]

	idI := h.getIdentifier(h.data[i])
	idJ := h.getIdentifier(h.data[j])
	h.indices[idI] = i
	h.indices[idJ] = j
}

func (h *Heap[T, I]) up(index int) {
	for index > 0 {
		parentIndex := parent(index)
		if !h.less(h.data[index], h.data[parentIndex]) {
			break
		}

		h.swap(index, parentIndex)
		index = parentIndex
	}
}

func (h *Heap[T, I]) down(index int) {
	minIndex := index
	size := len(h.data)

	for {
		leftIndex := leftChild(index)
		rightIndex := rightChild(index)

		if leftIndex < size && h.less(h.data[leftIndex], h.data[minIndex]) {
			minIndex = leftIndex
		}

		if rightIndex < size && h.less(h.data[rightIndex], h.data[minIndex]) {
			minIndex = rightIndex
		}

		if minIndex == index {
			break
		}

		h.swap(index, minIndex)
		index = minIndex
	}
}

func parent(i int) int {
	return (i - 1) / 2
}

func leftChild(i int) int {
	return 2*i + 1
}

func rightChild(i int) int {
	return 2*i + 2
}
