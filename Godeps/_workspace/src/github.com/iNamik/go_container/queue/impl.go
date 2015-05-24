package queue

type queue struct {
	data   []interface{}
	head   int
	tail   int
	length int
}

// Queue::Add
func (q *queue) Add(item interface{}) {
	if q.length == cap(q.data) {
		newData := make([]interface{}, cap(q.data)*2)
		for i := 0; i < q.length; i++ {
			newData[i] = q.data[q.head+i]
		}
		q.data = newData
		q.head = 0
		q.tail = q.length
	}
	q.data[q.tail] = item
	q.tail++
	q.length++
}

// Queue::Remove
func (q *queue) Remove() interface{} {
	if q.length == 0 {
		panic("underflow")
	}
	item := q.data[q.head]
	q.data[q.head] = nil
	q.head++
	q.length--
	if q.length == 0 {
		q.head = 0
		q.tail = 0
	}
	return item
}

// Queue::Peek
func (q *queue) Peek(n int) interface{} {
	if q.length == 0 || n < 0 || n > q.length {
		panic("range error")
	}
	return q.data[q.head+n]
}

// Queue::Clear
func (q *queue) Clear() {
	for i := q.head; i < q.tail; i++ {
		q.data[i] = nil
	}
	q.head = 0

	q.tail = 0

	q.length = 0
}

// Queue::Len
func (q *queue) Len() int {
	return q.length
}

// Queue::Cap
func (q *queue) Cap() int {
	return cap(q.data)
}

// Queue::Empty
func (q *queue) Empty() bool {
	return q.length == 0
}

// Queue::AtCapacity
func (q *queue) AtCapactiy() bool {
	return q.length == cap(q.data)
}
