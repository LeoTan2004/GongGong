package executor

import "sync"

// WorkerPool 定义一个协程池
type WorkerPool struct {
	tasks     chan func()    // 任务队列
	wg        sync.WaitGroup // 用于等待所有任务完成
	workerNum int            // 协程数量
	startup   bool
}

// NewWorkerPool 创建一个新的协程池
func NewWorkerPool(workerNum int) *WorkerPool {
	return &WorkerPool{
		tasks:     make(chan func()),
		workerNum: workerNum,
	}
}

// worker 是每个协程的具体执行逻辑
func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task()       // 执行任务
		wp.wg.Done() // 标记任务完成
	}
}

// Run 启动协程池，开启指定数量的协程
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workerNum; i++ {
		go wp.worker()
	}
}

// Submit 提交任务到任务队列
func (wp *WorkerPool) Submit(task func()) {
	if !wp.startup {
		wp.Run()
		wp.startup = true
	}
	wp.wg.Add(1) // 增加一个任务
	wp.tasks <- task

}

// Wait 等待所有任务完成
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()    // 等待任务完成
	close(wp.tasks) // 关闭任务队列
}
