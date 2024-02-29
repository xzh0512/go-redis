package consistenthash

import (
	"hash/crc32"
	"sort"
)

// HashFunc defines function to generate hash code
type HashFunc func(data []byte) uint32

// NodeMap stores nodes and you can pick node from NodeMap
// 一致性hash处理器，所有节点一致性哈希的信息
type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int          // sorted 各个节点的哈希值
	nodehashMap map[int]string //根据上面的哈希值找节点a,b,c..
}

// NewNodeMap creates a new NodeMap
func NewNodeMap(fn HashFunc) *NodeMap {
	//初始化
	m := &NodeMap{
		hashFunc:    fn,
		nodehashMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE //给一个初始的哈希函数
	}
	return m
}

// IsEmpty returns if there is no node in NodeMap
func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashs) == 0
}

// AddNode add the given nodes into consistent hash circle
// 把一个新的节点加入到一个环中
func (m *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashs = append(m.nodeHashs, hash)
		m.nodehashMap[hash] = key
	}
	sort.Ints(m.nodeHashs)
}

// PickNode gets the closest item in the hash to the provided key.
func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))

	// Binary search for appropriate replica.
	idx := sort.Search(len(m.nodeHashs), func(i int) bool {
		return m.nodeHashs[i] >= hash
	})

	// Means we have cycled back to the first replica.
	if idx == len(m.nodeHashs) {
		idx = 0
	}

	return m.nodehashMap[m.nodeHashs[idx]]
}
