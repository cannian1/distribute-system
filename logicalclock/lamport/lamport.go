package lamport

import "sync"

type LogicalClock struct {
	mu   sync.Mutex // 保证并发安全的互斥锁
	time int64      // 当前逻辑时间
}

// NewLogicalClock 初始化 lamport 时钟
func NewLogicalClock(initialTime int64) *LogicalClock {
	return &LogicalClock{
		time: initialTime,
	}
}

// Increment 处理本地事件并返回新的时间
func (lc *LogicalClock) Increment() int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.time++
	return lc.time
}

// Update 处理接收事件并返回调整后的时间
func (lc *LogicalClock) Update(receivedTime int64) int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	//  Lamport算法：time = max(local, received) + 1
	lc.time = max(lc.time, receivedTime) + 1

	return lc.time
}

// Now 获取当前逻辑时间（不修改时钟状态）
func (lc *LogicalClock) Now() int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	return lc.time
}

// Advance 手动推进时钟（仅在特殊情况下使用）
func (lc *LogicalClock) Advance(newTime int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if newTime > lc.time {
		lc.time = newTime
	}
}
