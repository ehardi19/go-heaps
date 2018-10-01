// Package pairing implements a Pairing heap Data structure
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Pairing_heap
package pairing

import (
	"go-heaps"
)

// Compile time check
var _ go_heaps.Interface = (*PairHeap)(nil)

// PairHeap represents a Pairing Heap.
// The zero value for PairHeap Root is an empty Heap.
type PairHeap struct {
	Comparator go_heaps.Comparator // Key comparator
	Root       *PairHeapNode
}

// PairHeapNode contains the current Value and the list if the sub-heaps
type PairHeapNode struct {
	// The current Top of the heap
	Value    interface{}
	// List of children PairHeapNodes all containing values less than the Top of the heap
	children []*PairHeapNode
	// A reference to the parent Heap Node
	parent   *PairHeapNode
}

// Init initializes or clears the PairHeap
func (p *PairHeap) Init(c go_heaps.Comparator) *PairHeap {
	p.Root = &PairHeapNode{}
	p.Comparator = c
	return p
}

// New returns an initialized PairHeap with the provided Comparator.
func New(c go_heaps.Comparator) *PairHeap { return new(PairHeap).Init(c) }

// NewWithIntComparator instantiates an PairHeap with the IntComparator, i.e. keys are of type int.
func NewWithIntComparator() *PairHeap {
	return new(PairHeap).Init(go_heaps.IntComparator)
}

// NewWithStringComparator instantiates an PairHeap with the StringComparator, i.e. keys are of type string.
func NewWithStringComparator() *PairHeap {
	return new(PairHeap).Init(go_heaps.StringComparator)
}

// IsEmpty returns true if PairHeap p is empty.
// The complexity is O(1).
func (p *PairHeap) IsEmpty() bool {
	return p.Root.Value == nil
}

// Resets the current PairHeap
func (p *PairHeap) Clear() {
	p.Root = &PairHeapNode{}
}

// Find the smallest item in the priority queue.
// The complexity is O(1).
func (p *PairHeap) FindMin() interface{} {
	if p.IsEmpty() {
		return nil
	}
	return p.Root.Value
}

// Inserts the value to the PairHeap and returns the PairHeapNode
// The complexity is O(1).
func (p *PairHeap) Insert(v interface{}) interface{}  {
	n := PairHeapNode{Value: v}
	merge(&p.Root, &n, p.Comparator)
	return v
}

// DeleteMin removes the top most value from the PairHeap and returns the combined PairHeap
// The complexity is O(log n) amortized.
func (p *PairHeap) DeleteMin() interface{} {
	if p.IsEmpty() {
		return nil
	}
	result := mergePairs(&p.Root, p.Root.children, p.Comparator)
	return result.Value
}

// Visits root and all of its children recursively invoking the callback function
func (p *PairHeap) Do(cb func(v interface{}))  {
	if p.IsEmpty() {
		return
	}
	// Call root first
	cb(p.Root)
	// Then continue to children
	visitChildren(p.Root.children, cb)
}

// Exhausting search of the element that matches v. Returns it as a PairHeapNode
// The complexity is O(n) amortized.
func (p *PairHeap) Find(v interface{}) *PairHeapNode {
	if p.IsEmpty() {
		return nil
	}

	if p.Comparator(p.Root.Value, v) == 0 {
		return p.Root
	} else {
		return p.findInChildren(p.Root.children, v)
	}
}

func (p *PairHeap) findInChildren(children []*PairHeapNode, v interface{}) *PairHeapNode  {
	if len(children) == 0 {
		return nil
	}
	var node *PairHeapNode
loop:
	for _, heapNode := range p.Root.children {
		cmp := p.Comparator(heapNode.Value, v)
		switch {
		case cmp == 0: // found
			node = heapNode
			break loop
		case cmp > 0:
			// current node has a higher value that item.
			// search in child nodes as they have lower value
			node = p.findInChildren(heapNode.children, v)
		case cmp < 0:
			// Continue as we know that the child nodes have even lower value
			continue
		}
	}
	return node
}

func visitChildren(children []*PairHeapNode, cb func(v interface{}))  {
	if len(children) == 0 {
		return
	}
	for _, heapNode := range children {
		cb(heapNode)
		visitChildren(heapNode.children, cb)
	}
}

func merge(first **PairHeapNode, second *PairHeapNode, c go_heaps.Comparator) *PairHeapNode {
	q := *first
	if q.Value == nil { // Case when root is empty
		*first = second
		return *first
	}

	cmp := c(q.Value, second.Value)
	if cmp < 0 {
		// put 'second' as the first child of 'first' and update the parent
		q.children = append([]*PairHeapNode{second}, q.children...)
		second.parent = *first
		return *first
	} else {
		// put 'first' as the first child of 'second' and update the parent
		second.children = append([]*PairHeapNode{q}, second.children...)
		q.parent = second
		*first = second
		return second
	}
}

// Merges heaps together
func mergePairs(root **PairHeapNode, heaps []*PairHeapNode, c go_heaps.Comparator) *PairHeapNode {
	q := *root
	if len(heaps) == 0 {
		*root = &PairHeapNode{
			parent: nil,
		}
		return q
	}
	if len(heaps) == 1 {
		*root = heaps[0]
		heaps[0].parent = nil
		return q
	}
	var merged *PairHeapNode
	for { // iteratively merge heaps
		if len(heaps) == 0 { break }
		if len(heaps) == 1 {
			// merge odd one out
			merged = merge(&merged, heaps[0], c)
			break
		}
		merged = merge(&heaps[0], heaps[1], c)
		heaps = heaps[2:]
	}
	*root = merged
	merged.parent = nil

	return q
}