package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64                    //最大内存限制
	nbytes   int64                    //当前已使用内存
	ll       *list.List               //双向链表
	cache    map[string]*list.Element //键值对-v是链表节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) //当被移除元素时的回调
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		//列表删除
		c.ll.Remove(ele)
		//kv删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//更新已使用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//是否回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	//更新或添加之后可能需要淘汰，不能直接return
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		//当前已使用 += 当前v新值大小 - v旧值大小
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { //缓存不存在
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		//当前已使用 += 当前k大小 + v大小
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	//添加缓存后若内存超过限制则淘汰
	//for 是没问题的，可能会 remove 多次，添加一条大的键值对，可能需要淘汰掉多个键值对，直到 nbytes < maxBytes
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
