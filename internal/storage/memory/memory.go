package memory

import (
	"container/heap"
	"github.com/harshabangi/url-shortener/internal/storage/shared"
	"sync"
)

func New() *memoryStore {
	return &memoryStore{
		urls:       make(map[string]string),
		domainFreq: make(map[string]int64),
		mutex:      sync.RWMutex{},
	}
}

type memoryStore struct {
	urls       map[string]string // Maps short codes to original URLs
	domainFreq map[string]int64  // Maps domain names to their frequencies
	mutex      sync.RWMutex      // Mutex for thread-safe access
}

func (m *memoryStore) SaveURL(key, originalURL string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.urls[key]; ok {
		return shared.ErrCollision
	}
	m.urls[key] = originalURL
	return nil
}

func (m *memoryStore) GetOriginalURL(key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if value, ok := m.urls[key]; ok {
		return value, nil
	}
	return "", shared.ErrNotFound
}

func (m *memoryStore) RecordDomainFrequency(domainName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.domainFreq[domainName]++
	return nil
}

func (m *memoryStore) GetTopNDomainsByFrequency(n int) ([]shared.DomainFrequency, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if len(m.domainFreq) <= 0 {
		return nil, nil
	}
	return getTopNDomainsByFrequency(m.domainFreq, n), nil
}

type domainFrequencyHeap []shared.DomainFrequency

func (h domainFrequencyHeap) Len() int           { return len(h) }
func (h domainFrequencyHeap) Less(i, j int) bool { return h[i].Frequency < h[j].Frequency }
func (h domainFrequencyHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *domainFrequencyHeap) Push(x interface{}) {
	*h = append(*h, x.(shared.DomainFrequency))
}

func (h *domainFrequencyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func getTopNDomainsByFrequency(domainFrequencyMap map[string]int64, n int) []shared.DomainFrequency {

	pq := make(domainFrequencyHeap, 0)
	heap.Init(&pq)

	for domain, frequency := range domainFrequencyMap {
		heap.Push(&pq, shared.DomainFrequency{Domain: domain, Frequency: frequency})
		if pq.Len() > n {
			heap.Pop(&pq)
		}
	}

	result := make([]shared.DomainFrequency, n)
	for i := n - 1; i >= 0; i-- {
		domainFrequency := heap.Pop(&pq).(shared.DomainFrequency)
		result[i] = domainFrequency
	}

	return result
}
