/**
 * @Author: dingQingHui
 * @Description:
 * @File: script_test
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:34
 */

package redis_extend

import (
	"testing"
	"time"
)

var testScript = NewRedisScript("testScript", `
	local data = ARGV[1]
	redis.call('SET', KEYS[1],  data)
	return 0
`)

func TestClient_Script(t *testing.T) {
	GetManager().Add(0, &Config{
		[]string{"192.168.1.25:6379"}, "", 100,
	})
	r, err := GetManager().GetById(0).Script(testScript, []string{"111111"}, "1")
	DebugLog(r, err)
}

func TestMsgQueue_Pop(t *testing.T) {
	GetManager().Add(0, &Config{
		[]string{"192.168.1.25:6379"}, "", 100,
	})

	mq := NewMsgQueue("testMsgQueue")
	mq.Push(GetManager().GetById(0), "222222")
	println(mq.PopWithTimeout(GetManager().GetById(0), time.Second*10))
}
