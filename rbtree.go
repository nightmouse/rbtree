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
	Value  int64
	left   *Node
	right  *Node
	parent *Node
}

// creates a new leaf node
func newLeafNode(parent *Node, value int64) *Node {
	return &Node{color: red,
		Value:  value,
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
			if n.Value == v {
				break
			} else if v < n.Value {
				if n.left == nil {
					n.left = newLeafNode(n, v)
					t.insertCase1(n.left)
					break
				}
				n = n.left
			} else if v > n.Value {
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
	if gp == nil {
		return
	}

	if n == n.parent.right && n.parent == gp.left { // n is a right child of a left child
		t.rotateLeft(n.parent)
		n = n.left
	} else if n == n.parent.left && n.parent == gp.right { // n is a left child of a right child
		t.rotateRight(n.parent)
		n = n.right
	}
	t.insertCase5(n)
}

func (t *RBTree) rotateLeft(n *Node) {
	rchild := n.right
	if n == t.root {
		t.root = rchild
	} else {
		if n == n.parent.left {
			n.parent.left = rchild
		} else if n == n.parent.right {
		    n.parent.right = rchild
        }
	}

	n.right = rchild.left
	rchild.parent = n.parent
    if rchild.left != nil {
        rchild.left.parent = n
    }
	rchild.left = n
	n.parent = rchild
}

func (t *RBTree) rotateRight(n *Node) {
	lchild := n.left

	if n == t.root {
		t.root = lchild
	} else {
		if n == n.parent.left {
			n.parent.left = lchild
		} else if n == n.parent.right {
			n.parent.right = lchild
		}
	}

    n.left = lchild.right
	lchild.parent = n.parent
    if lchild.right != nil {
        lchild.right.parent = n
    }
	lchild.right = n
	n.parent = lchild
}

func (t *RBTree) insertCase5(n *Node) {
	gp := grandparent(n)
	n.parent.color = black
	gp.color = red

	if n == n.parent.left {
		t.rotateRight(gp)
	} else if n == n.parent.right {
		t.rotateLeft(gp)
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


// a function to test the tree invariant, useful for debugging
func (t *RBTree) testInvariant(last int64) {
    fn := func(n *Node) {
        if n.parent == nil { return }
        if n == n.parent.left && !(n.Value < n.parent.Value) {
            panic(fmt.Sprintf("left child's value is greater than the parent value after insert %d", last))
        } else if n == n.parent.right && !(n.Value > n.parent.Value) {
            panic(fmt.Sprintf("right child's value is less than the parent value after insert %d", last))
        } else if n != n.parent.right && n != n.parent.left {
            panic(fmt.Sprintf("n != n.parent.right && n != n.parennt.left after node: %v parent: %v", n, n.parent))
        }
    }
    t.Do(TraverseBreadthFirst(fn))
}

// Applies fn to the root node.
func (t *RBTree) Do(fn func(*Node)) {
	fn(t.root)
}

// Returns true if val is found in the tree
func (t *RBTree) Find(val int64) bool {
	found := false
	f := func(n *Node) {
		for n != nil {
			if n.Value == val {
                found = true
				break
			} else if val < n.Value {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	t.Do(f)
	return found
}

// Creates a string by visting the nodes in order, with the format of:
// (value color) (value, color) (value, color)...
func (t *RBTree) String() string {
	buffer := &bytes.Buffer{}
	fn := func(n *Node) {
		format := "(%d %v %v %v %s)"
		if n == t.root {
			format = "(%d %v %v %v %s)*"
		}
		leftVal := "nil"
		rightVal := "nil"
        parentVal := "nil"
		if n.left != nil {
			leftVal = strconv.Itoa(int(n.left.Value))
		}
		if n.right != nil {
			rightVal = strconv.Itoa(int(n.right.Value))
		}
		if n.parent != nil {
			parentVal = strconv.Itoa(int(n.parent.Value))
		}
		buffer.WriteString(fmt.Sprintf(format, n.Value, leftVal, rightVal, parentVal, n.color))
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
		ch <- n.Value
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
		newTree.Insert(n.Value)
	}
	t.Do(TraverseBreadthFirst(fn))
	return newTree
}

// Slice returns an in-order slice of all the values in the tree.
func (t *RBTree) Slice() []int64 {
	slice := make([]int64, 0, t.Size())
	fn := func(n *Node) {
		slice = append(slice, n.Value)
	}
	t.Do(TraverseInOrder(fn))
	return slice
}
