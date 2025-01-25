package vector

import (
	"sync"
)

type Relation int

const (
	Before     Relation = iota // 当前时钟发生在另一个时钟之前
	After                      // 当前时钟发生在另一个时钟之后
	Concurrent                 // 两个时钟并发发生
	Equal                      // 两个时钟完全相同
)

// VectorClock 向量时钟
type VectorClock struct {
	nodeID     string            // 当前节点标识符
	timestamps map[string]uint64 // 所有节点的计数器状态
	mu         sync.RWMutex      // 读写锁保证线程安全
}

func NewVectorClock(nodeID string) *VectorClock {
	return &VectorClock{
		nodeID:     nodeID,
		timestamps: make(map[string]uint64),
	}
}

// Increment 增加当前节点的计数器值
func (vc *VectorClock) Increment() {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.timestamps[vc.nodeID]++
}

// Merge 合并另一个向量时钟的状态（取最大值）
func (vc *VectorClock) Merge(other *VectorClock) {
	otherCopy := other.cloneTimestamps()

	vc.mu.Lock()
	defer vc.mu.Unlock()

	for node, ts := range otherCopy {
		if current, ok := vc.timestamps[node]; !ok || ts > current {
			vc.timestamps[node] = ts
		}
	}
}

// Compare 比较两个向量时钟的关系
func (vc *VectorClock) Compare(other *VectorClock) Relation {
	vcCopy := vc.cloneTimestamps()
	otherCopy := other.cloneTimestamps()

	allLessOrEqual := true
	allGreaterOrEqual := true
	foundDifference := false

	// 创建所有节点的并集
	allNodes := make(map[string]struct{})
	for node := range vcCopy {
		allNodes[node] = struct{}{}
	}
	for node := range otherCopy {
		allNodes[node] = struct{}{}
	}

	// 比较每个节点的计数器值
	for node := range allNodes {
		v := vcCopy[node]
		o := otherCopy[node]

		if v < o {
			allGreaterOrEqual = false
		} else if v > o {
			allLessOrEqual = false
		}

		if v != o {
			foundDifference = true
		}
	}

	if !foundDifference {
		return Equal
	}
	if allLessOrEqual {
		return Before
	}
	if allGreaterOrEqual {
		return After
	}
	return Concurrent
}

// ToMap 返回当前向量时钟的副本（用于序列化）
func (vc *VectorClock) ToMap() map[string]uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	m := make(map[string]uint64, len(vc.timestamps))
	for k, v := range vc.timestamps {
		m[k] = v
	}
	return m
}

// FromMap 从map加载向量时钟状态（用于反序列化）
func (vc *VectorClock) FromMap(data map[string]uint64) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	vc.timestamps = make(map[string]uint64, len(data))
	for k, v := range data {
		vc.timestamps[k] = v
	}
}

// cloneTimestamps 复制当前时间戳状态（内部使用）
func (vc *VectorClock) cloneTimestamps() map[string]uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	m := make(map[string]uint64, len(vc.timestamps))
	for k, v := range vc.timestamps {
		m[k] = v
	}
	return m
}
