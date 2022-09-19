/**
 * @Author: dingQingHui
 * @Description:
 * @File: common
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:21
 */

package redis_extend

import "github.com/go-redis/redis"

// RedisNil redis空数据
func RedisNil(err error) bool {
	return err == redis.Nil
}

// RedisError  redis报错
func RedisError(err error) bool {
	if err == redis.Nil {
		return false
	}
	return err != nil
}
