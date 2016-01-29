package rbtree

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type color bool

func (c color) String() string {
	if c == black {
		return "black"
	}
	return "red"
}

const black color = false
const red color = true

// A red-black tree has the following properties:
// Property 1: every node is red or black
// Property 2: all leaf (nil) nodes are black
// Property 3: if a leaf node is red, the all it's children must be black
// Property 4: every path from a node to a leaf descendent has the same number of black nodes
type RBTree struct {
	nodeCount int
	root      *Node
}

type Node struct {
	color  color
	value  int64
	left   *Node
	right  *Node
	parent *Node
}

// creates a new leaf node
func newLeafNode(parent *Node, value int64) *Node {
	return &Node{color: red,
		value:  value,
		left:   nil,
		right:  nil,
		parent: parent}
}

func grandparent(n *Node) *Node {
	if n.parent != nil && n.parent.parent != nil {
		return n.parent.parent
	}
	return nil
}

func uncle(n *Node) *Node {
	gp := grandparent(n)
	if gp == nil {
		return nil
	}

	if gp.left == n.parent {
		return gp.right
	}
	return gp.left
}

// Create a new, empty red-black tree.
func New() *RBTree {
	return &RBTree{}
}

// Returns the number of nodes in the tree.
func (t *RBTree) Size() int {
	return t.nodeCount
}

// Calculates the tree height.
func (t *RBTree) Height() int64 {
	return 2 * int64(math.Log2(float64(t.nodeCount+1)))
}

func (t *RBTree) Insert(values ...int64) {
	for _, v := range values {
		t.nodeCount += 1
		if t.root == nil { // special case - nil root node
			t.root = newLeafNode(nil, v)
			t.root.color = black // Property 5: the root node is always black
			continue
		}

		n := t.root
		for {
			if n.value == v {
				break
			} else if v < n.value {
				if n.left == nil {
					n.left = newLeafNode(n, v)
					t.insertCase1(n.left)
					break
				}
				n = n.left
			} else if v > n.value {
				if n.right == nil {
					n.right = newLeafNode(n, v)
					t.insertCase1(n.right)
					break
				}
				n = n.right
			}
		} // end for loop
	} // end for loop
}

func (t *RBTree) rotateRight_case4(n *Node) {
	lchild := n.left
	if n == t.root {
		t.root = lchild
		t.root.parent = nil
	}
	n.left = lchild.right
	lchild.right = n
	lchild.parent = n.parent
	n.parent = lchild
}

func (t *RBTree) rotateLeft_case4(n *Node) {
	rchild := n.right
	if n == t.root {
		t.root = rchild
		t.root.parent = nil
	}
	n.right = rchild.left
	rchild.left = n
	rchild.parent = n.parent
	n.parent = rchild
}

func (t *RBTree) rotateLeft_case5(n *Node) {
	rchild := n.right
	rchild.parent = n.parent
	if n == t.root {
		t.root = rchild
	} else {
		rchild.parent.right = rchild
	}
	n.right = rchild.left
	rchild.left = n
	n.parent = rchild
}

func (t *RBTree) rotateRight_case5(n *Node) {
	lchild := n.left
	lchild.parent = n.parent
	if n == t.root {
		t.root = lchild
	} else {
		lchild.parent.left = lchild
	}
	n.left = lchild.right
	lchild.right = n
	n.parent = lchild
}

func (t *RBTree) insertCase1(n *Node) {
	if n.parent == nil {
		n.color = black
	} else {
		t.insertCase2(n)
	}
}

func (t *RBTree) insertCase2(n *Node) {
	if n.parent.color == red {
		t.insertCase3(n)
	}
}

func (t *RBTree) insertCase3(n *Node) {
	u := uncle(n)
	if u != nil && u.color == red {
		n.parent.color = black
		u.color = black
		gp := grandparent(n)
		gp.color = red
		t.insertCase1(gp)
	} else {
		t.insertCase4(n)
	}
}

func (t *RBTree) insertCase4(n *Node) {
	gp := grandparent(n)
	if n == n.parent.right && n.parent == gp.left { // n is a right child of a left child
		t.rotateLeft_case4(n)
		n = n.left
	} else if n == n.parent.left && n.parent == gp.right { // n is a left child of a right child
		t.rotateRight_case4(n)
		n = n.right
	}
	t.insertCase5(n)
}

func (t *RBTree) insertCase5(n *Node) {
	gp := grandparent(n)
	if gp == nil {
		return
	}
	n.parent.color = black
	gp.color = red
	if n == n.parent.left {
		t.rotateRight_case5(gp)
	} else if n == n.parent.right  {
		t.rotateLeft_case5(gp)
	}
}

// Create an annonymous function suitable for use with the Do method.
// The anonymous function applys fn over the nodes in depth first preorder.
func TraversePreOrder(fn func(*Node)) func(*Node) {
	var traverse func(*Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		fn(n)
		traverse(n.left)
		traverse(n.right)
	}
	return traverse
}

// Create an annonymous function suitable for use with the Do method.
// The anonymous function applys fn over the nodes in depth first in-order.
func TraverseInOrder(fn func(*Node)) func(*Node) {
	var traverse func(*Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		traverse(n.left)
		fn(n)
		traverse(n.right)
	}
	return traverse
}

// Create an annonymous function suitable for use with the Do method.
// The anonymous function applys fn over the nodes in depth first post order.
func TraversePostOrder(fn func(*Node)) func(*Node) {
	var traverse func(*Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		traverse(n.left)
		traverse(n.right)
		fn(n)
	}
	return traverse
}

// TraverseBreadthFirst creates an annonymous function suitable for use with the Do method.
// The anonymous function applys fn over the nodes in breadth first order.
func TraverseBreadthFirst(fn func(*Node)) func(*Node) {
	var innerFunc func(...*Node)
	innerFunc = func(nodes ...*Node) {
		children := make([]*Node, 0, 2*len(nodes))
		for _, n := range nodes {
			fn(n)
			if n.left != nil {
				children = append(children, n.left)
			}
			if n.right != nil {
				children = append(children, n.right)
			}
		}
		if len(children) > 0 {
			innerFunc(children...)
		}
	}

	traverse := func(n *Node) {
		innerFunc(n)
	}
	return traverse
}

// Applies fn to the root node.
func (t *RBTree) Do(fn func(*Node)) {
	fn(t.root)
}

// Creates a string by visting the nodes in order, with the format of:
// (value color) (value, color) (value, color)...
func (t *RBTree) String() string {
	buffer := &bytes.Buffer{}
	fn := func(n *Node) {
		format := "(%d %v %v %s)"
		if n == t.root {
			format = "(%d %v %v %s)*"
		}
		leftVal := "nil"
		rightVal := "nil"
		if n.left != nil {
			leftVal = strconv.Itoa(int(n.left.value))
		}
		if n.right != nil {
			rightVal = strconv.Itoa(int(n.right.value))
		}
		buffer.WriteString(fmt.Sprintf(format, n.value, leftVal, rightVal, n.color))
	}
	t.Do(TraverseInOrder(fn))
	return buffer.String()
}

// Return a channel suitable for use with range
func (t *RBTree) Iterate() <-chan int64 {
	ch := make(chan int64)
	if t.root == nil {
		close(ch)
		return ch
	}

	count := 0

	fn := func(n *Node) {
		count++
		ch <- n.value
		if count == t.nodeCount {
			close(ch)
		}
	}

	go t.Do(TraverseInOrder(fn))
	return ch
}

// Clone returns a deep copy of the current tree
func (t *RBTree) Clone() *RBTree {
	newTree := New()
	fn := func(n *Node) {
		newTree.Insert(n.value)
	}
	t.Do(TraverseBreadthFirst(fn))
	return newTree
}

// Slice returns an in-order slice of all the values in the tree.
func (t *RBTree) Slice() []int64 {
	slice := make([]int64, 0, t.Size())
	fn := func(n *Node) {
		slice = append(slice, n.value)
	}
	t.Do(TraverseInOrder(fn))
	return slice
}
