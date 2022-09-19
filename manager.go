/**
 * @Author: dingQingHui
 * @Description:
 * @File: Manager
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:11
 */

package redis_extend

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	once sync.Once
	mgr  *Manager
)

func GetManager() *Manager {
	once.Do(func() {
		mgr = new(Manager)
	})
	return mgr
}

type Manager struct {
	dbs sync.Map
}

func (m *Manager) GetById(id int) *Client {
	v, ok := m.dbs.Load(id)
	if !ok {
		return nil
	}
	c, ok := v.(*Client)
	if !ok {
		return nil
	}
	return c
}

func (m *Manager) Add(id int, conf *Config) {
	if c := m.GetById(id); c != nil {
		DebugLog("redis already exist", id)
		return
	}

	re := &Client{
		UniversalClient: redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    conf.Addrs,
			Password: conf.Passwd,
			PoolSize: conf.PoolSize,
		}),
		conf: conf,
	}

	m.dbs.Store(id, re)
}
