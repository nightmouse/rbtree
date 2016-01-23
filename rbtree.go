package rbtree

import (
	"bytes"
	"fmt"
	"math"
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

type RBTree struct {
	nodeCount uint64
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
func (t RBTree) Size() uint64 {
	return t.nodeCount
}

// Calculates the tree height.
func (t RBTree) Height() int64 {
	return 2 * int64(math.Log2(float64(t.nodeCount+1)))
}

func (t *RBTree) Insert(values ...int64) {
	for _, v := range values {
		t.nodeCount += 1
		if t.root == nil { // special case - nil root node
			t.root = newLeafNode(nil, v)
			t.root.color = black
			continue
		}

		n := t.root
		for {
			if n.value == v {
				break
			} else if v < n.value {
				if n.left == nil {
					n.left = newLeafNode(n, v)
					//rebalance(n.left)
					break
				}
				n = n.left
			} else if v > n.value {
				if n.right == nil {
					n.right = newLeafNode(n, v)
					//rebalance(n.right)
					break
				}
				n = n.right
			}
		} // end for loop
	} // end for loop
}

func (t *RBTree) rotateRight(n *Node) {
	lchild := n.left
	if n == t.root {
		t.root = lchild
	}
	n.left = lchild.right
	lchild.right = n
	lchild.parent = n.parent
	n.parent = lchild
}

func (t *RBTree) rotateLeft(n *Node) {
	rchild := n.right
	if n == t.root {
		t.root = rchild
	}
	n.right = rchild.left
	rchild.left = n
	rchild.parent = n.parent
	n.parent = rchild
}

// Property 1: every node is red or black
// Property 2: all leaf nodes are black
// Property 3: if a leaf node is red, the all it's children must be black
// Property 4: every path from a node to a leaf descendent has the same number of black nodes
// Property 5: the root node is always black

//func (t *RBTree) rebalance(n *Node) {
//    if n.parent == nil {
//      return
//    }
//
//    gp : = grandparent(n)
//    u := uncle(n)
//    if gp == nil || u == nil {
//      return
//    }
//
//    // case 1: uncle is red
//    // swap colors of uncle, parent, and grand parent
//    if u.color == red {
//        n.parent.color = !n.parent.color
//        u.color =  !u.color
//        gp.color = !gp.color
//        return
//    }
//
//    // case 2 & 3: uncle is black
//
//    // case 2
//    if n.parent.left == n && gp.right == n.parent {   // n is a left child of a right child
//      rotateRight(n.parent)
//    } else if n.parent.right == n && gp.left == n.parent { // n is a right child or a left child
//      rotateLeft(n.parent)
//    }
//
//    // get the grandparent and uncle since they've (maybe) changed in the rotate
//    gp : = grandparent(n)
//    u := uncle(n)
//
//    // case 3
//    if n.parent. {
//      rotateRight(gp)
//    } else if {
//      rotateLeft(gp)
//    }
//
//}

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

// Create an annonymous function suitable for use with the Do method.
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
		buffer.WriteString(fmt.Sprintf("(%d, %s)", n.value, n.color))
	}
	t.Do(TraverseInOrder(fn))
	return buffer.String()
}

// Return a channel suitable for use with range
func (t *RBTree) Iterate() <-chan int64 {
	ch := make(chan int64)
	count := uint64(0)

	fn := func(n *Node) {
		ch <- n.value
		count++
		if count == t.nodeCount {
			close(ch)
		}
	}

	go t.Do(TraverseInOrder(fn))
	return ch
}

func (t *RBTree) Clone() *RBTree {
	newTree := New()
	fn := func(n *Node) {
		newTree.Insert(n.value)
	}
	t.Do(TraverseBreadthFirst(fn))
	return newTree
}

// Returns an in-order slice of all the values in the tree.
func (t *RBTree) Slice() []int64 {
	slice := make([]int64, 0, t.Size())
	fn := func(n *Node) {
		slice = append(slice, n.value)
	}
	t.Do(TraverseInOrder(fn))
	return slice
}
