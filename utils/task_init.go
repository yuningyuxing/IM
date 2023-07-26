package utils

import "time"

type TimerFunc func(interface{}) bool

// 定时函数  用于定时执行指定的函数
// delay表示首次延迟时间 表示多久后执行函数
// tick 表示间隔时间 表示函数的间隔时间
// TimerFunc表示定时执行的函数
// param表示传入定时执行函数的参数
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	//启动一个协程 以免阻塞主函数
	go func() {
		//如果定时执行的函数为空 直接返回 不执行定时函数
		if fun == nil {
			return
		}
		//创建一个定时器 并设置首次延迟时间
		t := time.NewTimer(delay)
		//开始无限循环 并且等待定时器的C通道接受值
		for {
			select {
			//定时器的C通道会在定时器到期后接受一个值 此时执行定时任务
			case <-t.C:
				//执行定时函数 并传入参数param 如果返回false 不再执行定时函数
				if fun(param) == false {
					return
				}
				//重置定时器 以时间间隔的定时执行
				t.Reset(tick)
			}
		}
	}()
}
