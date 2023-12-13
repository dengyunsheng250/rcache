package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 一致性hash算法将key映射到2^32的空间中，形成一个环
// 计算key的哈希值，放置在环上，顺时针寻找到的第一个节点，就是应选取的节点/机器。
// 解决数据倾斜问题
// 引入虚拟节点的概念, 一个真实节点对应多个虚拟节点
// 第一步，计算虚拟节点的 Hash 值，放置在环上。
// 第二步，计算 key 的 Hash 值，在环上顺时针寻找到应选取的虚拟节点，例如是 peer2-1，那么就对应真实节点 peer2。
// 代价很小

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash  // hash函数
	replicas int   // 虚拟节点倍数
	keys     []int //
	hashMap  map[int]string
}

func New(hash Hash, replicas int) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	// 对所有虚拟节点进行排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	if idx == len(m.keys) {
		idx = 0
	}
	// 一致性哈希是个环状结构
	return m.hashMap[m.keys[idx]]
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
