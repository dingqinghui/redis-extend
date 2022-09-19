/**
 * @Author: dingQingHui
 * @Description:redis实现异步消息队列
 * @File: mq
 * @Version: 1.0.0
 * @Date: 2022/9/19 15:15
 */

package redis_extend

import "time"

type MsgQueue struct {
	key string
}

func NewMsgQueue(key string) *MsgQueue {
	return &MsgQueue{
		key: key,
	}
}

func (m *MsgQueue) Push(client *Client, values ...string) error {
	_, err := client.LPush(client.Context(), m.key, values).Result()
	return err
}

func (m *MsgQueue) PopWithTimeout(client *Client, timeout time.Duration) (string, error) {
	r, err := client.BRPop(client.Context(), timeout, m.key).Result()
	if err != nil {
		if RedisNil(err) {
			return "", nil
		}
		return "", err
	}
	return r[1], nil
}
