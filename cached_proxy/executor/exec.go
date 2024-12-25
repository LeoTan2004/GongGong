package executor

type Executor interface {
	Run()               // 开始执行
	Submit(task func()) // 提交任务
	Wait()              // 等待执行完成
}
