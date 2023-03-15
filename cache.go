package fate

type Cache struct {
	names []FirstName
}

func (c *Cache) Put(name FirstName) {
	c.names = append(c.names, name)
}

func (c *Cache) Len() int {
	return len(c.names)
}

func (c *Cache) Get(sta, end int) []FirstName {
	if sta >= c.Len() || sta < 0 {
		return nil
	}
	if end >= c.Len() || end < 0 {
		return nil
	}
	if sta > end {
		sta, end = end, sta
	}
	return c.names[sta:end]
}

func (c *Cache) GetOne(idx int) (FirstName, bool) {
	if idx >= c.Len() || idx < 0 {
		return FirstName{}, false
	}
	return c.names[idx], true
}
