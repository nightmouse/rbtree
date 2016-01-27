package rbtree

import (
	"reflect"
	"testing"
)

var testData = [][]int64{
	{},
	{5},
	{1, 2},
	{1, 2, 3},
	{1, 2, 3, 4},
	{1, 2, 3, 4, 5},
	{1, 2, 3, 4, 5, 6},
}

var reverseTestData = [][]int64{
	{2, 1},
	{3, 2, 1},
	{4, 3, 2, 1},
	{5, 4, 3, 2, 1},
	{6, 5, 4, 3, 2, 1},
}

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
	for _, td := range testData {
		tree := New()
		tree.Insert(td...)
	}
}

func TestString(t *testing.T) {
	var expected = []string{
		"",
		"(5 nil nil black)*",
		"(1 nil 2 black)*(2 nil nil red)",
		"(1 nil nil red)(2 1 3 black)*(3 nil nil red)",
		"(1 nil nil black)(2 1 3 black)*(3 nil 4 black)(4 nil nil red)",
		"(1 nil nil black)(2 1 4 black)*(3 nil nil red)(4 3 5 black)(5 nil nil red)"}

	for i, td := range testData[:len(expected)] {
		tree := New()
		tree.Insert(td...)
		str := tree.String()
		if expected[i] != str {
			t.Errorf("\nexpected %v\nactual   %v", expected[i], str)
		}
	}

	var reverseExpected = []string{
		"(1 nil nil red)(2 1 nil black)*",
		"(1 nil nil red)(2 1 3 black)*(3 nil nil red)",
		"(1 nil nil red)(2 1 nil black)(3 2 4 black)*(4 nil nil black)"}

	for i, td := range reverseTestData[:len(reverseExpected)] {
		tree := New()
		tree.Insert(td...)
		str := tree.String()
		if reverseExpected[i] != str {
			t.Errorf("\nexpected %v\nactual   %v", reverseExpected[i], str)
		}
	}
}

func TestSize(t *testing.T) {
	for _, td := range testData {
		tree := New()
		tree.Insert(td...)
		if tree.Size() != len(td) {
			t.Error("Insert failed on input with ", len(td), " elements")
		}
	}
}

func TestFail(t *testing.T) {
	expected := "(1 nil nil red)(2 1 3 black)*(3 nil nil red)"
	data := []int64{3, 2, 1}
	tree := New()
	tree.Insert(data...)
	str := tree.String()
	if expected != str {
		t.Errorf("\nexpected %v\nactual   %v", expected, str)
	}
}

func TestRange(t *testing.T) {
	for _, td := range testData {
		tree := New()
		tree.Insert(td...)
		if tree.Size() != len(td) {
			t.Error("Insert failed on input with ", len(td), " elements")
		}

		actual := make([]int64, 0, len(td))
		for n := range tree.Iterate() {
			actual = append(actual, n)
		}

		if !reflect.DeepEqual(td, actual) {
			t.Error("Iterate failed on ", td)
		}
	}
}

func TestBasicRotateRight(t *testing.T) {
	expected := []int64{4, 5, 2, 1, 3}
	tree := New()
	tree.Insert(expected...)
	t.Log("before: ", tree.String())
	tree.rotateRight_case4(tree.root)
	t.Log("after: ", tree.String())
}

func TestBasicRotateLeft(t *testing.T) {
	expected := []int64{2, 1, 4, 3, 5}
	tree := New()
	tree.Insert(expected...)
	t.Log("before: ", tree.String())
	tree.rotateLeft_case4(tree.root)
	t.Log("after: ", tree.String())
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
	t.Log(slice)
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
