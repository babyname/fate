package fate

// IteratorFunc ...
type IteratorFunc func(v interface{}) error

type iterator struct {
	data  []interface{}
	index int
}

// newIterator ...
func newIterator() *iterator {
	return &iterator{
		data:  nil,
		index: 0,
	}
}

//HasNext check next
func (i *iterator) HasNext() bool {
	return i.index < len(i.data)
}

//Next get next
func (i *iterator) Next() interface{} {
	defer func() {
		i.index++
	}()
	if i.index < len(i.data) {
		return i.data[i.index]
	}

	return nil
}

//Reset reset index
func (i *iterator) Reset() {
	i.index = 0
}

//Add add radical
func (i *iterator) Add(v interface{}) {
	i.data = append(i.data, v)
}

//Size iterator data size
func (i *iterator) Size() int {
	return len(i.data)
}

//Iterator an default iterator
func (i *iterator) Iterator(f IteratorFunc) error {
	i.Reset()
	for i.HasNext() {
		if err := f(i.Next()); err != nil {
			return err
		}
	}
	return nil
}
