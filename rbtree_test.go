package rbtree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tree := New()
	if tree.String() != "" {
		t.Error("empty tree somehow made a non-empty string")
	}
}

func TestSimpleColor(t *testing.T) {
	// root node should always be black
	tree := New()
	tree.Insert(2)
	if tree.root.color != black {
		t.Error("After 1 insertion: the root node should be black")
	}

	tree.Insert(3)
	tree.Insert(1)
	// the root node should still be black, and the child node should be red
	if tree.root.color != black {
		t.Error("After 3 insertions: the root node should be black")
	}

	if tree.root.right.color != red {
		t.Error("the right leaf node should be red, but is not")
	}

	if tree.root.left.color != red {
		t.Error("The left leaf node should be read but is not")
	}
}

func TestInsert(t *testing.T) {
	tree := New()
	tree.Insert(5, 4, 7, 3, 2, 6, 8, 9)
	fmt.Println(tree)
}

func TestSize(t *testing.T) {
	{
		tree := New()
		if tree.Size() != 0 {
			t.Error("Size() failed on an empty trree")
		}
	}

	{
		tree := New()
		values := []int64{1}
		tree.Insert(values...)
		if tree.Size() != len(values) {
			t.Error("Size() failed after one insertion")
		}
	}

	{
		tree := New()
		values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		tree.Insert(values...)
		if tree.Size() != len(values) {
			t.Error("Size() failed after %d insertion", len(values))
		}
	}

}

func TestRange(t *testing.T) {
	{
		tree := New()
		for _ = range tree.Iterate() {
			t.Error("an empty tree should never iterate")
		}
	}

	{
		tree := New()
		expected := []int64{1}
		tree.Insert(expected...)
		actual := 0
		for _ = range tree.Iterate() {
			actual++
		}
		if actual != tree.Size() {
			t.Error("iteration failed after iteration of one item")
		}
	}

	{
		tree := New()
		expected := []int64{5, 4, 7, 3, 2, 6, 8, 9}
		fmt.Println("expected length: ", len(expected))

		tree.Insert(expected...)
		fmt.Println(tree)
		for _ = range tree.Iterate() {
			//fmt.Println(v)
		}
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
	expected := []int64{1, 2, 3, 4, 5}
	tree := New()
	tree.Insert(values...)
	slice := tree.Slice()
	fmt.Println(slice)
	if !reflect.DeepEqual(slice, expected) {
		t.Error("Expected ", expected, " got ", slice)
	}
}

func BenchmarkInsert1k(b *testing.B) {
	tree := New()
	for i := 0; i != b.N; i++ {
		tree.Insert(int64(i))
	}
}
