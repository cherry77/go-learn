package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/**
题目：并发任务处理器
编写一个 Go 程序实现以下功能：
1. 并发执行多个任务：
    * 使用sync.WaitGroup管理并发执行的 goroutine
    * 每个任务模拟不同的处理时间（随机生成）
2. 任务结果收集：
    * 使用通道收集任务执行结果（成功/失败）
    * 每个任务结果应包含：任务ID、执行时间、是否成功、错误信息（如果有）
3. 进度显示：
    * 实时显示每个任务的执行进度（百分比）
    * 使用单独的进度通道
4. 统计功能：
    * 所有任务完成后，显示：
        * 总执行时间
        * 成功任务数
        * 失败任务数
        * 平均任务执行时间
        * 最快/最慢任务信息
5. 额外要求：
    * 实现并发数控制（最多同时运行 n 个任务）
    * 任务失败后可以重试（最多重试 m 次）
*/

// Task 表示一个任务
type Task struct {
	ID         int
	RetryCount int
	MaxRetries int
}

// TaskResult 表示任务执行结果
type TaskResult struct {
	TaskID     int
	Duration   time.Duration
	IsSuccess  bool
	Error      string
	RetryCount int
}

// ProgressUpdate 表示进度更新
type ProgressUpdate struct {
	TaskID    int
	Progress  float64
	IsRunning bool
}

func main() {
	// 配置参数
	totalTasks := 10
	maxConcurrent := 3
	maxRetries := 2

	// 创建任务
	tasks := make([]Task, totalTasks)
	for i := 0; i < totalTasks; i++ {
		tasks[i] = Task{
			ID:         i + 1,
			RetryCount: 0,
			MaxRetries: maxRetries,
		}
	}

	// 创建通道
	resultChan := make(chan TaskResult, totalTasks)
	progressChan := make(chan ProgressUpdate, totalTasks*100) // 假设每个任务最多发送100个进度更新
	done := make(chan struct{})

	// 启动结果收集器
	var results []TaskResult
	go func() {
		for result := range resultChan {
			results = append(results, result)
		}
		close(done)
	}()

	// 启动进度显示器
	go func() {
		progressMap := make(map[int]float64)
		for update := range progressChan {
			progressMap[update.TaskID] = update.Progress
			fmt.Printf("\r")
			for id, p := range progressMap {
				if p < 100 {
					fmt.Printf("任务 %d: %.1f%% | ", id, p)
				}
			}
		}
		fmt.Println("\n所有任务完成!")
	}()

	// 记录开始时间
	startTime := time.Now()

	// 使用 WaitGroup 管理并发
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrent) // 控制并发数

	// 启动所有任务
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			semaphore <- struct{}{} // 获取信号量
			executeTask(t, resultChan, progressChan)
			<-semaphore // 释放信号量
		}(task)
	}

	// 等待所有任务完成
	wg.Wait()
	close(resultChan)
	close(progressChan)

	// 等待结果收集完成
	<-done

	// 计算统计信息
	calculateStatistics(results, startTime)
}

// executeTask 执行单个任务
func executeTask(task Task, resultChan chan<- TaskResult, progressChan chan<- ProgressUpdate) {
	var result TaskResult
	var err error
	var duration time.Duration

	for task.RetryCount <= task.MaxRetries {
		// 模拟任务执行
		duration, err = simulateTask(task.ID, progressChan)
		if err == nil {
			// 任务成功
			result = TaskResult{
				TaskID:     task.ID,
				Duration:   duration,
				IsSuccess:  true,
				Error:      "",
				RetryCount: task.RetryCount,
			}
			resultChan <- result
			return
		}

		// 任务失败
		task.RetryCount++
		if task.RetryCount <= task.MaxRetries {
			fmt.Printf("\n任务 %d 失败，准备第 %d 次重试...\n", task.ID, task.RetryCount)
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1)) // 重试前等待
		}
	}

	// 重试次数用完仍失败
	result = TaskResult{
		TaskID:     task.ID,
		Duration:   duration,
		IsSuccess:  false,
		Error:      err.Error(),
		RetryCount: task.RetryCount - 1, // 显示实际重试次数
	}
	resultChan <- result
}

// simulateTask 模拟任务执行
func simulateTask(taskID int, progressChan chan<- ProgressUpdate) (time.Duration, error) {
	startTime := time.Now()
	duration := time.Duration(rand.Intn(5)+1) * time.Second // 随机持续时间 1-5秒
	failureRate := 0.3                                      // 30%失败率

	// 发送进度更新
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(startTime)
			progress := float64(elapsed) / float64(duration) * 100
			if progress > 100 {
				progress = 100
			}
			progressChan <- ProgressUpdate{
				TaskID:    taskID,
				Progress:  progress,
				IsRunning: true,
			}
			if progress >= 100 {
				// 随机决定任务成功或失败
				if rand.Float64() > failureRate {
					return duration, nil
				}
				return duration, fmt.Errorf("任务 %d 模拟失败", taskID)
			}
		}
	}
}

// calculateStatistics 计算并显示统计信息
func calculateStatistics(results []TaskResult, startTime time.Time) {
	var (
		successCount  int
		failCount     int
		totalDuration time.Duration
		minDuration   time.Duration
		maxDuration   time.Duration
	)

	if len(results) == 0 {
		fmt.Println("没有任务结果可统计")
		return
	}

	minDuration = results[0].Duration
	maxDuration = results[0].Duration

	for _, result := range results {
		totalDuration += result.Duration

		if result.Duration < minDuration {
			minDuration = result.Duration
		}
		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}

		if result.IsSuccess {
			successCount++
		} else {
			failCount++
		}
	}

	avgDuration := totalDuration / time.Duration(len(results))
	totalTime := time.Since(startTime)

	fmt.Println("\n===== 统计信息 =====")
	fmt.Printf("总执行时间: %v\n", totalTime.Round(time.Millisecond))
	fmt.Printf("成功任务数: %d\n", successCount)
	fmt.Printf("失败任务数: %d\n", failCount)
	fmt.Printf("平均任务执行时间: %v\n", avgDuration.Round(time.Millisecond))
	fmt.Printf("最快任务: %v (任务ID: ", minDuration.Round(time.Millisecond))
	for _, result := range results {
		if result.Duration == minDuration {
			fmt.Printf("%d ", result.TaskID)
		}
	}
	fmt.Println(")")

	fmt.Printf("最慢任务: %v (任务ID: ", maxDuration.Round(time.Millisecond))
	for _, result := range results {
		if result.Duration == maxDuration {
			fmt.Printf("%d ", result.TaskID)
		}
	}
	fmt.Println(")")

	// 打印失败任务详情
	if failCount > 0 {
		fmt.Println("\n失败任务详情:")
		for _, result := range results {
			if !result.IsSuccess {
				fmt.Printf("任务 %d: 重试 %d 次, 错误: %s\n",
					result.TaskID, result.RetryCount, result.Error)
			}
		}
	}
}
