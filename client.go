/**
 * @Author: dingQingHui
 * @Description:
 * @File: redis
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:05
 */

package redis_extend

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"strings"
)

var (
	ErrDBErr      = errors.New("数据库错误")
	ErrDBDataType = errors.New("数据库数据类型错误")
)

type Client struct {
	redis.UniversalClient
	conf *Config
}

// ScriptStr 脚本执行
func (r *Client) ScriptStr(cmd int, keys []string, args ...interface{}) (string, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		ErrorLog("redis script failed", err)
		return "", ErrDBErr
	}
	_, ok := data.(int64)
	if ok {
		return "", ErrDBErr
	}
	if data == nil {
		return "", nil
	}
	str, ok := data.(string)
	if !ok {
		return "", ErrDBDataType
	}
	return str, nil
}

// Script 脚本处理
func (r *Client) Script(cmd int, keys []string, args ...interface{}) (interface{}, error) {
	var err error = ErrDBErr
	var re interface{}
	// 腾讯云Redis必须填一个Key
	// keys = append(keys, "bug{tag}")
	hashStr, ok := scriptHashMap.Load(cmd)
	if ok {
		re, err = r.EvalSha(r.Context(), hashStr.(string), keys, args...).Result()
	}
	if RedisError(err) {
		scriptStr, ok1 := scriptMap.Load(cmd)
		if !ok1 {
			ErrorLog("redis script error cmd not found", cmd)
			return nil, ErrDBErr
		}
		cmdStr, _ := scriptCommitMap.Load(cmd)
		if strings.HasPrefix(err.Error(), "NOSCRIPT ") || err == ErrDBErr {
			DebugLog("try reload redis script", cmdStr.(string))
			hashStr, err = r.ScriptLoad(r.Context(), scriptStr.(string)).Result()
			if err != nil {
				ErrorLog("redis script load", err, cmdStr.(string))
				return nil, ErrDBErr
			}
			scriptHashMap.Store(cmd, hashStr.(string))
			re, err = r.EvalSha(r.Context(), hashStr.(string), keys, args...).Result()
			if !RedisError(err) {
				return re, nil
			}
		}
		ErrorLog("redis script error", err, cmdStr.(string))
		return nil, ErrDBErr
	}

	return re, nil
}
