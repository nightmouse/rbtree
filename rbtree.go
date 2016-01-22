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
	root      *node
}

type node struct {
	color  color
	value  int64
	left   *node
	right  *node
	parent *node
}

// creates a new leaf node
func newLeafNode(parent *node, value int64) *node {
	return &node{color: red,
		value:  value,
		left:   nil,
		right:  nil,
		parent: parent}
}

func grandparent(n *node) *node {
	if n.parent != nil && n.parent.parent != nil {
		return n.parent.parent
	}
	return nil
}

func uncle(n *node) *node {
	gp := grandparent(n)
	if gp == nil {
		return nil
	}

	if gp.left == n.parent {
		return gp.right
	}
	return gp.left
}

func New() *RBTree {
	return &RBTree{}
}

func (t RBTree) Size() uint64 {
	return t.nodeCount
}

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

func (t *RBTree) rotateRight(n *node) {
    lchild := n.left
    if n == t.root { 
        t.root = lchild 
    }
    n.left = lchild.right 
    lchild.right = n
    lchild.parent = n.parent
    n.parent = lchild
}

func (t *RBTree)rotateLeft(n *node) {
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

//func (t *RBTree) rebalance(n *node) {
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

// (value color) (value, color) (value, color)
func (t *RBTree) String() string {
	buffer := &bytes.Buffer{}
	fn := func(n *node) {
		buffer.WriteString(fmt.Sprintf("(%d, %s)", n.value, n.color))
	}

	t.Do(fn)
	return buffer.String()
}

// TODO: add other traversal methods
// applies fn to each node in pre-order traversal
func (t *RBTree) Do(fn func(*node)) {
	var preorderTraverse func(n *node)
	preorderTraverse = func(n *node) {
		if n == nil {
			return
		}

		fn(n)

		preorderTraverse(n.left)
		preorderTraverse(n.right)
	}
	preorderTraverse(t.root)
}

func (t *RBTree) Iterate() <- chan int64 {
	ch := make(chan int64)
	count := uint64(0)

	fn := func(n *node) {
		ch <- n.value
		count++
		if count == t.nodeCount {
			close(ch)
		}
	}

	go t.Do(fn)
	return ch
}

// TODO: there's got to be a more efficient way to do this
func (t *RBTree) Clone() *RBTree { 
    newTree := New()
    fn := func(n *node) { 
        newTree.Insert(n.value)
    }
    t.Do(fn)
    return newTree
}

// returns a slice of all the values in the tree, in pre-order traversal
// TODO: return a sorted array when multiple traversals are supported
func (t *RBTree) Slice() []int64 { 
	slice := make([]int64, 0, t.Size())
	fn := func(n *node) {
		slice = append(slice, n.value)	
	}
	t.Do(fn)
	return slice
}
