/**
 * @Author: dingQingHui
 * @Description:
 * @File: config
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:06
 */

package redis_extend

type Config struct {
	Addrs    []string
	Passwd   string
	PoolSize int
}
