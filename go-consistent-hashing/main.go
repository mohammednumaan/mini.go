package main

// this implementation was referenced from:
// https://github.com/arpitbbhayani/consistent-hashing/blob/master/consistent-hashing.ipynb

// the article is at:
// https://arpitbhayani.me/blogs/consistent-hashing/

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"sort"
)

type StorageNode struct {
	host string
	name string
}

type ConsistentHashing struct {
	nodePositions []int
	nodes         []StorageNode
	totalSlots    int
}

func NewConsistentHashing() *ConsistentHashing {
	return &ConsistentHashing{
		nodePositions: []int{},
		nodes:         []StorageNode{},
		totalSlots:    50,
	}
}

func hashFn(value string, totalSlots int) int {
	h := sha256.Sum256([]byte(value))

	hashInt := new(big.Int).SetBytes(h[:])
	mod := new(big.Int).Mod(hashInt, big.NewInt(int64(totalSlots)))

	return int(mod.Int64())
}

func (ch *ConsistentHashing) addNode(node StorageNode) (int, error) {
	if len(ch.nodePositions) == ch.totalSlots {
		return -1, errors.New("hash space is full")
	}

	nodePosition := hashFn(node.host, ch.totalSlots)
	idx := sort.SearchInts(ch.nodePositions, nodePosition)

	if idx < len(ch.nodePositions) && ch.nodePositions[idx] == nodePosition {
		return -1, errors.New("collision occurred")
	}

	ch.nodePositions = append(ch.nodePositions, 0)
	ch.nodes = append(ch.nodes, StorageNode{})

	copy(ch.nodePositions[idx+1:], ch.nodePositions[idx:])
	copy(ch.nodes[idx+1:], ch.nodes[idx:])

	ch.nodePositions[idx] = nodePosition
	ch.nodes[idx] = node

	return nodePosition, nil
}

func (ch *ConsistentHashing) removeNode(node StorageNode) (int, error) {
	if len(ch.nodePositions) == 0 {
		return -1, errors.New("hash space is empty")
	}

	nodePosition := hashFn(node.host, ch.totalSlots)
	idx := sort.SearchInts(ch.nodePositions, nodePosition)

	if idx >= len(ch.nodePositions) || ch.nodePositions[idx] != nodePosition {
		return -1, errors.New("node does not exist")
	}

	ch.nodePositions = append(ch.nodePositions[:idx], ch.nodePositions[idx+1:]...)
	ch.nodes = append(ch.nodes[:idx], ch.nodes[idx+1:]...)

	return nodePosition, nil
}

func (ch *ConsistentHashing) assign(item string) (StorageNode, error) {
	if len(ch.nodePositions) == 0 {
		return StorageNode{}, errors.New("no nodes in ring")
	}

	itemPosition := hashFn(item, ch.totalSlots)

	idx := sort.Search(len(ch.nodePositions), func(i int) bool {
		return ch.nodePositions[i] > itemPosition
	})

	idx = idx % len(ch.nodePositions)
	return ch.nodes[idx], nil
}

func main() {
	ch := NewConsistentHashing()

	_, _ = ch.addNode(StorageNode{host: "10.0.0.1", name: "A"})
	_, _ = ch.addNode(StorageNode{host: "10.0.0.2", name: "B"})
	_, _ = ch.addNode(StorageNode{host: "10.0.0.3", name: "C"})

	node1, _ := ch.assign("random_photo.png")
	node2, _ := ch.assign("random_audio.mp3")
	node3, _ := ch.assign("random_video.mov")

	fmt.Printf("Item 'random_photo.png' is assigned to node: %s\n", node1.name)
	fmt.Printf("Item 'random_audio.mp3' is assigned to node: %s\n", node2.name)
	fmt.Printf("Item 'random_video.mov' is assigned to node: %s\n", node3.name)
}
