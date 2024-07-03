package cache

type Cache struct {
	Items map[string]map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		Items: make(map[string]map[string]interface{}),
	}
}
func (c *Cache) Get(topic string, key string) interface{} {
	return c.Items[topic][key]
}
func (c *Cache) Set(topic string, key string, item interface{}) {
	c.Items[topic][key] = item
}
