package utils

import "sync"

/////////////////////////////////////
//             Element             //
/////////////////////////////////////
// https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/container/list/list.go;l=48
type Element struct { // Linked List 의 Element
	next, prev *Element
	list       *List
	Value      any
}

func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

/////////////////////////////////////
//             List                //
/////////////////////////////////////
type List struct {
	root Element
	len  int
}

func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}
func (l *List) lazyInit() { // https://sonagi87174.tistory.com/9  // 임시로 초기화 하고 나중에 실재로 사용할떄 Init 을 하는 기능
	if l.root.next == nil {
		l.Init()
	}
}
func NewList() *List {
	return new(List).Init()
}
func (l *List) Len() int {
	return l.len
}
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}
func (l *List) insert(e, at *Element) *Element { // == insert(e *Element, at *Element) 랑 같은 의미
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at)
func (l *List) insertValue(v any, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}
func (l *List) remove(e *Element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}
func (l *List) move(e, at *Element) {
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
func (l *List) Remove(e *Element) any {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}
func (l *List) PushFront(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}
func (l *List) PushBack(v any) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}
func (l *List) InsertBefore(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark.prev)
}
func (l *List) InsertAfter(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark)
}

/////////////////////////////////////
//              Queue              //
/////////////////////////////////////

type Queue struct {
	v   *List
	len int
}

func (q *Queue) Push(T any) {
	q.v.PushBack(T)
	q.len++
}
func (q *Queue) Pop() any {
	front := q.v.Front()
	if front != nil {
		q.len--
		return q.v.Remove(front)
	}
	return nil
}

func NewQueue() *Queue {
	return &Queue{NewList(), 0}
}

type Stack struct {
}

/////////////////////////////////////
//              Room               //
/////////////////////////////////////

type Room struct {
	lock sync.Mutex
}
