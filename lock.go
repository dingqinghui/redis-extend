/**
 * @Author: dingQingHui
 * @Description:redis 实现不可重入分布式锁
 * @File: lock
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:56
 */

package redis_extend

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	DefaultRetry    = 5
	DefaultInterval = 100 * time.Millisecond
	DefaultExpire   = 5 * time.Second
)

var ErrFailed = errors.New("failed to acquire lock")

type LockConfig struct {
	//
	// Expiry
	// @Description:过期时间应大于业务执行时间,否则锁失效
	//
	Expiry        time.Duration
	Retry         int
	RetryInterval time.Duration
	SameLock      bool
}

var DefaultLockConfig = &LockConfig{
	Expiry:        DefaultExpire,
	RetryInterval: DefaultInterval,
	Retry:         DefaultRetry,
	SameLock:      false,
}

type Lock struct {
	lockKey   string
	lockValue string
	mutex     sync.Mutex
	rds       *Client
	*LockConfig
}

func NewLock(key string, cfg *LockConfig) *Lock {
	if cfg == nil {
		return nil
	}
	if len(key) <= 0 {
		return nil
	}
	return &Lock{
		lockKey:    key,
		LockConfig: cfg,
	}
}

func (rl *Lock) Lock() error {
	return rl.lock()
}

func (rl *Lock) lock() error {
	if rl.SameLock {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		rl.mutex.Lock()
		rl.lockValue = base64.StdEncoding.EncodeToString(b)
		rl.mutex.Unlock()
	} else {
		rl.lockValue = ""
	}

	for i := 0; i < rl.Retry; i++ {
		ok, err := rl.rds.SetNX(rl.rds.Context(), rl.lockKey, rl.lockValue, rl.Expiry).Result()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if i < rl.Retry-1 {
			time.Sleep(rl.RetryInterval)
		}
	}
	return ErrFailed
}

var (
	RedisScriptUnlock = NewRedisScript("RedisScriptAddMail", `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end`)
)

func (rl *Lock) UnLock() error {
	var err error
	if rl.SameLock {
		_, err = rl.rds.Script(RedisScriptUnlock, []string{rl.lockKey}, rl.lockValue)
		rl.mutex.Unlock()
		return err
	} else {
		_, err = rl.rds.Del(rl.rds.Context(), rl.lockKey).Result()

	}
	if err != nil {
		ErrorLog("redis unlock fail", rl.lockKey, rl.lockValue)
	}
	return err
}
