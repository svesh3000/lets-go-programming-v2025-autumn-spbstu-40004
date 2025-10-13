package MaxHeap

type MaxHeap []int

func (heap *MaxHeap) Push(val any) {
	value, ok := val.(int)
	if ok {
		*heap = append(*heap, value)
	}

	panic("MaxHeap.Push: value must be int")
}

func (heap *MaxHeap) Pop() any {
	len := heap.Len()
	if len == 0 {
		return nil
	}

	element := (*heap)[len-1]
	*heap = (*heap)[0 : len-1]

	return element
}

func (heap *MaxHeap) Len() int {
	return len(*heap)
}

func (heap *MaxHeap) Less(i, j int) bool {
	return (*heap)[i] > (*heap)[j]
}

func (heap *MaxHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}
