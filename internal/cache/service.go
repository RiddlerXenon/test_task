package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Item struct {
	Value string
	TTL   int64
}

type Cache struct {
	Items   map[string]Item
	Cleaner *cleaner
	Mutex   sync.Mutex
}

func New() (*Cache, error) {
	cln := &cleaner{
		closed: make(chan bool),
	}

	c := &Cache{
		Items:   make(map[string]Item),
		Cleaner: cln,
	}

	go cln.Run(c)
	runtime.SetFinalizer(c, closeCleaner)

	return c, nil
}

func (c *Cache) Add(key, value string, ttl time.Duration) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, exists := c.Items[key]; exists {
		return fmt.Errorf("key already exists")
	}
	c.Items[key] = Item{
		Value: value,
		TTL:   time.Now().Add(ttl).UnixNano(),
	}
	return nil
}

func (c *Cache) Set(key, val string, ttl time.Duration) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if item, exists := c.Items[key]; exists {
		item.Value = val
		item.TTL = time.Now().Add(ttl).UnixNano()
		return nil
	}
	return fmt.Errorf("key %s not found", key)
}

func (c *Cache) Get(key string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	item, exists := c.Items[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	if item.Value == "" {
		return "", fmt.Errorf("value for key %s is empty", key)
	}

	return item.Value, nil
}

func (c *Cache) Del(key string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, exists := c.Items[key]; exists {
		delete(c.Items, key)
		return nil
	}
	return fmt.Errorf("key not found")
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for key, item := range c.Items {
		if now >= item.TTL {
			delete(c.Items, key)
		}
	}
}
