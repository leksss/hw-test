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
	keys     map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	var wasInCache bool

	if len(c.items) == c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		c.removeLeastRecentlyUsed(lastItem)
	}

	if _, wasInCache = c.items[key]; wasInCache {
		c.queue.Remove(c.items[key])
	}

	item := c.queue.PushFront(value)
	c.items[key] = item
	c.keys[item] = key

	return wasInCache
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if elem, ok := c.items[key]; ok {
		c.queue.PushFront(c.items[key].Value)
		return elem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue.RemoveAll()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) removeLeastRecentlyUsed(item *ListItem) {
	if key, ok := c.keys[item]; ok {
		delete(c.items, key)
	}
}
