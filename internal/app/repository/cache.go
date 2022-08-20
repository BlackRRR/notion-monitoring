package repository

import (
	"sync"
)

type StatusCache struct {
	sync.Mutex
	TaskStatus map[string]string
}

type DescriptionCache struct {
	sync.Mutex
	TaskDescription map[string]string
}

func NewStatusCache() *StatusCache {
	return &StatusCache{
		Mutex:      sync.Mutex{},
		TaskStatus: make(map[string]string),
	}
}

func NewDescriptionCache() *DescriptionCache {
	return &DescriptionCache{
		Mutex:           sync.Mutex{},
		TaskDescription: make(map[string]string),
	}
}

func (c *StatusCache) Add(key, value string) map[string]string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.TaskStatus[key] = value

	return c.TaskStatus
}

func (c *StatusCache) Get(key string) string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	value := c.TaskStatus[key]

	return value
}

func (c *DescriptionCache) Add(key, value string) map[string]string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.TaskDescription[key] = value

	return c.TaskDescription
}

func (c *DescriptionCache) Get(key string) string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	value := c.TaskDescription[key]

	return value
}
