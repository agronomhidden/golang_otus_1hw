package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	queue    List
	capacity int
	sync.Mutex
	items map[Key]*cacheItem
}

func (e *lruCache) Set(key Key, value interface{}) bool {
	e.Lock()
	defer e.Unlock()
	if item, ok := e.items[key]; ok {
		item.value = value
		_ = e.queue.MoveToFront(item.parent) //кеш инкапсулирует лист, что гарантирует отсутствие кривых элементов из вне
		return true
	}
	e.items[key] = &cacheItem{key: key, value: value}
	qItem := e.queue.PushFront(e.items[key])
	e.items[key].parent = qItem
	e.normalizeCapacity()

	return false
}

func (e *lruCache) Get(key Key) (interface{}, bool) {
	e.Lock()
	defer e.Unlock()
	if item, ok := e.items[key]; ok {
		_ = e.queue.MoveToFront(item.parent)
		return item.value, true
	}
	return nil, false
}

func (e *lruCache) Clear() {
	e.items = make(map[Key]*cacheItem)
	e.queue = NewList()
}

func (e *lruCache) normalizeCapacity() {
	if e.queue.Len() > e.capacity {
		last := e.queue.Back()
		cItem := last.Value.(*cacheItem)
		delete(e.items, cItem.key)
		_ = e.queue.Remove(last)
	}
}

type cacheItem struct {
	key    Key
	value  interface{}
	parent *listItem //В идеале, List нужно отправить в отдельный пакет, т.к. listItem приватный, то нужно будет сделать публичный интерфейс для элемента, чтобы использовать его в типизации для кеша
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*cacheItem),
	}
}
