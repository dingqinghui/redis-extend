/**
 * @Author: dingQingHui
 * @Description:
 * @File: log
 * @Version: 1.0.0
 * @Date: 2022/9/19 11:14
 */

package redis_extend

func DebugLog(args ...interface{}) {
	println(args)
}

func ErrorLog(args ...interface{}) {
	println(args)
}
