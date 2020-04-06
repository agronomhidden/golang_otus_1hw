package hw04_lru_cache //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var mutex = &sync.Mutex{}
var counter uint
var ErrOuterItem = errors.New("list is not contain received item")

func getId() uint {
	mutex.Lock()
	defer mutex.Unlock()
	counter++
	return counter
}
// хочется сделать защиту от передачи в список элемента,
// который не принадлежит данному списку
// для этого вводим некий идентификатор элемента
// и ошибки если передали чужеродный элемент
type List interface {
	Len() int
    Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem) error
	MoveToFront(i *listItem) error
}

// думаю есть необходимость защитить next и prev, чтобы извне нельзя было переопределить соседа
// на элемент из совсем другого списка.
type listItem struct {
	Value interface{}
	next *listItem
	prev *listItem
	id uint
}

func (e listItem) Next() *listItem {
	return e.next
}

func (e listItem) Prev() *listItem {
	return e.prev
}

type list struct {
	first *listItem
	last *listItem
	index map[uint]struct{}
	sync.Mutex
	len int
}

func (e list) Len() int {
	return e.len
}

func (e *list) PushFront(v interface{}) *listItem {
	newItem := e.initNewItem(v)
	newItem.next = e.first
	if e.first == nil {
		e.last = newItem
	}
	e.first = newItem

	return newItem
}

func (e *list) PushBack(v interface{}) *listItem {
	newItem := e.initNewItem(v)
	newItem.prev = e.last
	if e.last == nil {
		e.first = newItem
	}
	e.last = newItem

	return newItem
}

func (e *list) Remove(item *listItem) error {
	if _, ok := e.index[item.id]; !ok {
		return ErrOuterItem
	}
	e.Lock()
	defer e.Unlock()

	e.len--
	item.id = 0

	if item.prev == nil && item.next == nil {
		e.first = nil
		e.last  = nil
	} else if item.prev == nil {
		e.first = item.next
		item.next = nil
	} else if item.next == nil {
		e.last = item.prev
		item.prev = nil
	} else {
		item.prev.next = item.next
		item.next.prev = item.prev
		item.next = nil
		item.prev = nil
	}
	return nil
}

func (e *list) MoveToFront(item *listItem) error {
	if _, ok := e.index[item.id]; !ok {
		return ErrOuterItem
	}
	if e.first == item {
		return nil
	}
	e.Lock()
	defer e.Unlock()

	item.prev.next = item.next
	item.next = e.first
	e.first = item

	return nil
}

func (e list) Front() *listItem {
	return e.first
}

func (e list) Back() *listItem {
	return e.last
}

func (e *list) initNewItem(v interface{}) *listItem {
	e.Lock()
	defer e.Unlock()

	newItem := &listItem{Value: v, id: getId()}
	e.index[newItem.id] = struct{}{}
	e.len++

	return newItem
}

func NewList() List {
	return &list{}
}
