package executor

import (
	"sync/atomic"
	"testing"
)

// 测试线程池基本运行情况
func TestWorkerPool(t *testing.T) {
	workerNum := 5
	taskCount := 20

	// 创建 WorkerPool 实例
	wp := NewWorkerPool(workerNum)

	// 启动 WorkerPool
	wp.Run()

	// 使用 atomic 计数器记录任务执行情况
	var executedTaskCount int32 = 0

	// 提交多个任务
	for i := 0; i < taskCount; i++ {
		wp.Submit(func() {
			atomic.AddInt32(&executedTaskCount, 1) // 原子增加任务计数
		})
	}

	// 等待所有任务完成
	wp.Wait()

	// 验证执行的任务数量是否与提交的任务数量一致
	if executedTaskCount != int32(taskCount) {
		t.Errorf("执行的任务数量错误：期望 %d，实际 %d", taskCount, executedTaskCount)
	}
}

func TestWorkerPoolRunMultipleTimes(t *testing.T) {
	workerNum := 3
	taskCount := 10

	// 创建 WorkerPool 实例
	wp := NewWorkerPool(workerNum)

	// 启动 WorkerPool
	wp.Run()

	// 使用 atomic 计数器记录任务执行情况
	var executedTaskCount int32 = 0

	// 提交多个任务
	for i := 0; i < taskCount; i++ {
		wp.Submit(func() {
			atomic.AddInt32(&executedTaskCount, 1)
		})
	}

	// 等待所有任务完成
	wp.Wait()

	// 验证执行的任务数量是否与提交的任务数量一致
	if executedTaskCount != int32(taskCount) {
		t.Errorf("执行的任务数量错误：期望 %d，实际 %d", taskCount, executedTaskCount)
	}
}
