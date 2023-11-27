package hash

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"slices"
	"strings"

	"github.com/zeromicro/go-zero/core/lang"
)

const (
	// TopWeight is the top weight that one entry might set.
	TopWeight = 100

	minReplicas = 2
	prime       = 16777619
)

type (
	// Func defines the hash method.
	Func func(data []byte) uint64

	// A ConsistentHash is a ring hash implementation.
	ConsistentHash struct {
		hashFunc Func
		replicas int
		keys     []uint64
		ring     map[uint64][]any
		nodes    map[string]string // String is the IP address
		lock     sync.RWMutex
	}
)

// NewConsistentHash returns a ConsistentHash.
func NewConsistentHash() *ConsistentHash {
	return NewCustomConsistentHash(minReplicas, Hash)
}

// NewCustomConsistentHash returns a ConsistentHash with given replicas and hash func.
func NewCustomConsistentHash(replicas int, fn Func) *ConsistentHash {
	if replicas < minReplicas {
		replicas = minReplicas
	}

	if fn == nil {
		fn = Hash
	}

	return &ConsistentHash{
		hashFunc: fn,
		replicas: replicas,
		ring:     make(map[uint64][]any),
		nodes:    make(map[string]string),
	}
}

// Add adds the node with the number of h.replicas,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) Add(node any, address string) {
	h.AddWithReplicas(node, h.replicas, address)
}

// AddWithReplicas adds the node with the number of replicas,
// replicas will be truncated to h.replicas if it's larger than h.replicas,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithReplicas(node any, replicas int, address string) {
	h.Remove(node)

	if replicas > h.replicas {
		replicas = h.replicas
	}

	nodeRepr := repr(node)
	h.lock.Lock()
	defer h.lock.Unlock()
	h.addNode(nodeRepr, address)

	for i := 0; i < replicas; i++ {
		hash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i)))
		h.keys = append(h.keys, hash)
		h.ring[hash] = append(h.ring[hash], node)
	}

	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})
}

// AddWithWeight adds the node with weight, the weight can be 1 to 100, indicates the percent,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithWeight(node any, weight int, address string) {
	// don't need to make sure weight not larger than TopWeight,
	// because AddWithReplicas makes sure replicas cannot be larger than h.replicas
	replicas := h.replicas * weight / TopWeight
	h.AddWithReplicas(node, replicas, address)
}

// Get returns the corresponding node from h base on the given v.
func (h *ConsistentHash) Get(v any) (any, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if len(h.ring) == 0 {
		return nil, false
	}

	hash := h.hashFunc([]byte(repr(v)))
	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	}) % len(h.keys)

	nodes := h.ring[h.keys[index]]

	switch len(nodes) {
		case 0:
			return nil, false
		case 1:
			return nodes[0], true
		default:
			innerIndex := h.hashFunc([]byte(innerRepr(v)))
			pos := int(innerIndex % uint64(len(nodes)))
			return nodes[pos], true
	}
}

// GetClosestNodes returns a slice of the numNodes closest nodes in the consistent hash ring
// to the provided value v. The function calculates the hash of v, finds the index of the
// first key greater than or equal to the hash, and retrieves the nodes associated with that
// key and the next two keys in the hash ring. If successful, it returns the slice of closest
// nodes and true; otherwise, it returns nil and false.
func (h *ConsistentHash) GetClosestNodes(v any, numNodes int) ([]any, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if len(h.ring) == 0 {
		return nil, false
	}

	hash := h.hashFunc([]byte(repr(v)))

	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	}) % len(h.keys)

	var closestNodes []any
	for i := 0; (len(closestNodes) < numNodes) && (len(closestNodes) < len(h.nodes)); i++ {
		keyIndex := (index + i) % len(h.keys)
		nodes := h.ring[h.keys[keyIndex]]

		for _, elem := range(nodes){
			if(!slices.Contains(closestNodes, elem)){
				closestNodes = append(closestNodes, elem)
			}
		}
	}

	if len(closestNodes) == 0 {
		return nil, false
	}

	return closestNodes, true
}

func (h *ConsistentHash) GetServerIP(server string) string {
	return h.nodes[server]
}

func (h *ConsistentHash) GetNumberOfKeys() int{
	return len(h.nodes)
}


func (h *ConsistentHash) GetServerName() string {
	var numbers []int

	for key := range h.nodes {

		numberStr := strings.TrimPrefix(key, "server ")
		number, err := strconv.Atoi(numberStr)

		if err == nil {
			numbers = append(numbers, number)
		}
	}

	sort.Ints(numbers)

	// Find the smallest available number
	smallestNumber := 0
	for _, num := range numbers {
		if num == smallestNumber {
			smallestNumber++
		} else {
			break
		}
	}

	return fmt.Sprintf("server %d", smallestNumber)
}


func (h *ConsistentHash) RemoveNodeByIP(address string) {

	for key, value := range(h.nodes){
		if(value == address){
			h.Remove(key)
			return
		}
	}

}

// Remove removes the given node from h.
func (h *ConsistentHash) Remove(node any) {
	nodeRepr := repr(node)

	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.containsNode(nodeRepr) {
		return
	}

	for i := 0; i < h.replicas; i++ {
		hash := h.hashFunc([]byte(nodeRepr + strconv.Itoa(i)))
		index := sort.Search(len(h.keys), func(i int) bool {
			return h.keys[i] >= hash
		})
		if index < len(h.keys) && h.keys[index] == hash {
			h.keys = append(h.keys[:index], h.keys[index+1:]...)
		}
		h.removeRingNode(hash, nodeRepr)
	}

	h.removeNode(nodeRepr)
}

// Get returns the map of nodes.
func (h *ConsistentHash) GetNodes() map[string]string{
	return h.nodes
}

// Get returns the map ring.
func (h *ConsistentHash) GetRing() map[uint64][]any{
	return h.ring
}

func (h *ConsistentHash) removeRingNode(hash uint64, nodeRepr string) {
	if nodes, ok := h.ring[hash]; ok {
		newNodes := nodes[:0]
		for _, x := range nodes {
			if repr(x) != nodeRepr {
				newNodes = append(newNodes, x)
			}
		}
		if len(newNodes) > 0 {
			h.ring[hash] = newNodes
		} else {
			delete(h.ring, hash)
		}
	}
}

func (h *ConsistentHash) addNode(nodeRepr string, address string) {
	h.nodes[nodeRepr] = address
}

func (h *ConsistentHash) containsNode(nodeRepr string) bool {
	_, ok := h.nodes[nodeRepr]
	return ok
}

func (h *ConsistentHash) removeNode(nodeRepr string) {
	delete(h.nodes, nodeRepr)
}

func innerRepr(node any) string {
	return fmt.Sprintf("%d:%v", prime, node)
}

func repr(node any) string {
	return lang.Repr(node)
}