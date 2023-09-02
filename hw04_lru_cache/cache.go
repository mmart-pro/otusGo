package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// пришлось ввести отдельную структуру для хранения в ListItem пары ключ-значение, чтобы можно было удалить элемент по ключу
type cacheElement struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, found := c.items[key]
	if found {
		// - если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
		item.Value = &cacheElement{key: key, value: value}

		c.queue.MoveToFront(item)
	} else {
		// - если элемента нет в словаре, то добавить в словарь и в начало очереди
		item = c.queue.PushFront(&cacheElement{key: key, value: value})
		c.items[key] = item

		if c.queue.Len() > c.capacity {
			last := c.queue.Back()
			c.queue.Remove(last)
			// тут нам нужен ключ последнего элемента в очереди для удаления из items
			k := last.Value.(*cacheElement).key
			delete(c.items, k)
		}
	}
	return found
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, found := c.items[key]
	// - если элемента нет в словаре, то вернуть nil и false
	if !found {
		return nil, false
	}
	// - если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
	c.queue.MoveToFront(item)

	return item.Value.(*cacheElement).value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
