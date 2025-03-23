package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type GenericNode[K comparable, V any] struct {
	Key   K
	Value V
	Left  *GenericNode[K, V]
	Right *GenericNode[K, V]
}

func (n *GenericNode[K, V]) ApplyAction(action func(K, V)) {
	if n == nil {
		return
	}

	if n.Left != nil {
		n.Left.ApplyAction(action)
	}
	action(n.Key, n.Value)
	if n.Right != nil {
		n.Right.ApplyAction(action)
	}
}

func (n *GenericNode[K, V]) GetLeftMostNode() *GenericNode[K, V] {
	if n == nil {
		return nil
	}

	node := n
	for node.Left != nil {
		node = node.Left
	}

	return node
}

type OrderedMap[K comparable, V any] struct {
	root *GenericNode[K, V]
	size int
	less func(a, b K) bool
}

func NewOrderedMap[K comparable, V any](less func(a, b K) bool) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		root: nil,
		size: 0,
		less: less,
	}
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if m == nil {
		var zeroValue V
		return zeroValue, false
	}

	node, ok := m.findNode(key)
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	return node.Value, true
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m == nil {
		return
	}

	m.root = m.insertNode(m.root, key, value)
	m.size++
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if m == nil || m.root == nil {
		return
	}

	m.root = m.removeNode(m.root, key)
	m.size--
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	if m == nil {
		return false
	}

	_, ok := m.findNode(key)

	return ok
}

func (m *OrderedMap[K, V]) Size() int {
	if m == nil || m.root == nil {
		return 0
	}

	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	if m == nil {
		return
	}

	m.root.ApplyAction(func(key K, value V) {
		action(key, value)
	})
}

func (m *OrderedMap[K, V]) findNode(key K) (*GenericNode[K, V], bool) {
	if m.root == nil {
		return nil, false
	}

	node := m.root
	for node != nil {
		if key == node.Key {
			return node, true
		}

		if m.less(key, node.Key) {
			node = node.Left
		} else {
			node = node.Right
		}
	}

	return nil, false
}

func (m *OrderedMap[K, V]) insertNode(node *GenericNode[K, V], key K, value V) *GenericNode[K, V] {
	if node == nil {
		return &GenericNode[K, V]{
			Key:   key,
			Value: value,
			Left:  nil,
			Right: nil,
		}
	}

	if m.less(key, node.Key) {
		node.Left = m.insertNode(node.Left, key, value)
	} else {
		node.Right = m.insertNode(node.Right, key, value)
	}

	return node
}

func (m *OrderedMap[K, V]) removeNode(node *GenericNode[K, V], key K) *GenericNode[K, V] {
	if node == nil {
		return nil
	}

	if key == node.Key {
		if node.Left == nil {
			return node.Right
		}

		if node.Right == nil {
			return node.Left
		}

		minNode := node.Right.GetLeftMostNode()

		node.Key = minNode.Key
		node.Value = minNode.Value
		node.Right = m.removeNode(node.Right, minNode.Key)

		return node
	}

	if m.less(key, node.Key) {
		node.Left = m.removeNode(node.Left, key)
	} else {
		node.Right = m.removeNode(node.Right, key)
	}

	return node
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](func(a, b int) bool {
		return a < b
	})
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
