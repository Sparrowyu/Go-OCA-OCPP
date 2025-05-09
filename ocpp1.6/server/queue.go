package server

import "sync"

type Queue interface {
	Push(v interface{})
	Pop() (interface{}, bool)
	Peek() (interface{}, bool)
	Len() int
	IsEmpty() bool
}

//minQueueLen must be power of 2   x % n == x & (n - 1)
const minQueueLen = 2 << 4

type lockQueue struct {
	buf               []interface{}
	head, tail, count int
	sync.RWMutex
}

func NewRequestQueue() *lockQueue {
	return &lockQueue{
		buf: make([]interface{}, minQueueLen),
	}
}

func NewEpollEventsQueue() *lockQueue {
	return &lockQueue{
		buf: make([]interface{}, minQueueLen),
	}
}

func (q *lockQueue) Len() int {
	q.RLock()
	defer q.RUnlock()
	return q.count
}

func (q *lockQueue) IsEmpty() bool {
	q.RLock()
	defer q.RUnlock()
	return q.count == 0
}

func (q *lockQueue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *lockQueue) Push(elem interface{}) {
	q.Lock()
	defer q.Unlock()
	if q.count == len(q.buf) {
		q.resize()
	}
	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

func (q *lockQueue) Peek() (interface{}, bool) {
	q.RLock()
	defer q.RUnlock()
	if q.count <= 0 {
		return nil, false
	}
	return q.buf[q.head], true
}

func (q *lockQueue) Pop() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.count <= 0 {
		return nil, false
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret, true
}
