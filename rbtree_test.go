package rbtree

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tree := New()
	if tree.String() != "" {
		t.Error("empty tree somehow made a non-empty string")
	}
}

func TestInsert(t *testing.T) {
	tree := New()
	tree.Insert(5, 4, 7, 3, 2, 6, 8, 9)
	fmt.Println(tree)
}

func TestRange(t *testing.T) {
    expected := []int64{5, 4, 7, 3, 2, 6, 8, 9}
	tree := New()
	tree.Insert(expected...)
    for v := range tree.Iterate() { 
        fmt.Println(v)
    }
}

func TestBasicRotateRight(t *testing.T) {
    expected := []int64{4, 5, 2, 1, 3}
	tree := New()
	tree.Insert(expected...)
    fmt.Println("before: ", tree.String())
    tree.rotateRight(tree.root)
    fmt.Println("after: ", tree.String())
}

func TestBasicRotateLeft(t *testing.T) {
    expected := []int64{2, 1, 4, 3, 5}
	tree := New()
	tree.Insert(expected...)
    fmt.Println("before: ", tree.String())
    tree.rotateLeft(tree.root)
    fmt.Println("after: ", tree.String())
}

func TestSmallClone(t *testing.T) {
    expected := []int64{4, 5, 2, 1, 3}
	tree := New()
	tree.Insert(expected...)
    newTree := tree.Clone()
    if tree.String() != newTree.String() { 
		t.Error("Expected ", tree.String(), " got ", newTree.String())
    }
}

func TestSmallSlice(t *testing.T) { 
    values := []int64{4, 5, 2, 1, 3}
    //expected:= []int64{4, 5, 2, 1, 3}
    tree := New()
    tree.Insert(values...)
    slice := tree.Slice()
    fmt.Println(slice)
}
