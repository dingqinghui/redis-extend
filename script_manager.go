/**
 * @Author: dingQingHui
 * @Description:
 * @File: script_manager
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:30
 */

package redis_extend

import (
	"sync"
	"sync/atomic"
)

var (
	scriptMap             = sync.Map{} // map[int]string{}
	scriptCommitMap       = sync.Map{} // map[int]string{}
	scriptHashMap         = sync.Map{} // map[int]string{}
	scriptIndex     int32 = 0
)

// NewRedisScript 新建一个Redis脚本
func NewRedisScript(commit, str string) int {
	cmd := int(atomic.AddInt32(&scriptIndex, 1))
	scriptMap.Store(cmd, str)
	scriptCommitMap.Store(cmd, commit)
	return cmd
}
