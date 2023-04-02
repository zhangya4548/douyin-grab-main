package cache

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

type Cache struct {
	cache *bigcache.BigCache
}

func NewCache() *Cache {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(72 * time.Hour))
	if err != nil {
		panic("chache初始化异常")
	}
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) SetDefaultCaChe() error {
	err := c.Set("LiveRoomUrl", "")
	if err != nil {
		return err
	}
	err = c.Set("WssUrl", "")
	if err != nil {
		return err
	}
	err = c.Set("Stop", "false")
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	// 获取缓存值
	value, err := c.cache.Get(key)
	if err != nil {
		return "", fmt.Errorf("获取缓存值%s出错：%s", key, err)
	}
	return string(value), nil
}

func (c *Cache) Del(key string) error {
	// 删除缓存值
	err := c.cache.Delete(key)
	if err != nil {
		return fmt.Errorf("删除缓存值%s出错：%s", key, err)
	}
	return nil
}

func (c *Cache) Set(key, val string) error {
	// 设置缓存值
	err := c.cache.Set(key, []byte(val))
	if err != nil {
		return fmt.Errorf("设置缓存值%s出错：%s", key, err)
	}
	return nil
}
