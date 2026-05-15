package session

// Element 双向链表的元素节点，包含指向前后节点的指针和存储的值。
type Element[T any] struct {
	next, prev *Element[T]
	list       *List[T]
	Value      T
}

// Next 返回当前元素的下一个元素，若已是末尾则返回 nil。
func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev 返回当前元素的上一个元素，若已是开头则返回 nil。
func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// List 泛型双向链表，使用哨兵根节点简化边界操作。
type List[T any] struct {
	root Element[T]
	len  int
}

// Init 初始化或重置链表，返回链表自身。
func (l *List[T]) Init() *List[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New 创建并初始化一个空的泛型双向链表。
func New[T any]() *List[T] { return new(List[T]).Init() }

// Len 返回链表中元素的个数。
func (l *List[T]) Len() int { return l.len }

// Front 返回链表的第一个元素，若链表为空则返回 nil。
func (l *List[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back 返回链表的最后一个元素，若链表为空则返回 nil。
func (l *List[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

func (l *List[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *List[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

func (l *List[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--
}

func (l *List[T]) move(e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove 从链表中移除指定元素，并返回其存储的值。
func (l *List[T]) Remove(e *Element[T]) T {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}

// PushFront 在链表头部插入一个新元素，并返回该元素。
func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack 在链表尾部插入一个新元素，并返回该元素。
func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore 在指定标记元素之前插入一个新元素，若标记不属于该链表则返回 nil。
func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark.prev)
}

// InsertAfter 在指定标记元素之后插入一个新元素，若标记不属于该链表则返回 nil。
func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark)
}

// MoveToFront 将指定元素移动到链表头部。
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	l.move(e, &l.root)
}

// MoveToBack 将指定元素移动到链表尾部。
func (l *List[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.move(e, l.root.prev)
}

// MoveBefore 将元素 e 移动到标记元素 mark 之前。
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter 将元素 e 移动到标记元素 mark 之后。
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList 将另一个链表的所有元素追加到当前链表尾部。
func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList 将另一个链表的所有元素逆序插入到当前链表头部。
func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
