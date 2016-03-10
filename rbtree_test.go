package rbtree

import (
	"math/rand"
	"reflect"
	"sort"
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
	{1, 2, 3, 4, 5, 6, 7},
	{1, 2, 3, 4, 5, 6, 8},
	{1, 2, 3, 4, 5, 6, 8, 9},
	{1, 2, 3, 4, 5, 6, 8, 9, 10},
	{4, 3, 5, 2, 1},
}

var reverseTestData = [][]int64{
	{2, 1},
	{3, 2, 1},
	{4, 3, 2, 1},
	{5, 4, 3, 2, 1},
	{6, 5, 4, 3, 2, 1},
	{7, 6, 5, 4, 3, 2, 1},
	{8, 6, 5, 4, 3, 2, 1},
	{9, 8, 6, 5, 4, 3, 2, 1},
	{10, 9, 8, 6, 5, 4, 3, 2, 1},
}

type Int64Slice []int64

func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s Int64Slice) Len() int {
	return len(s)
}
func (s Int64Slice) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

func scramble(seed int64, slice []int64) []int64 {
	rand.Seed(seed)
	size := len(slice)
	s := make([]int64, size)
	copy(s, slice)
	for i := 0; i != size; i++ {
		j := rand.Intn(size)
		s[j], s[i] = s[i], s[j]
	}
	return s
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
		"(5 nil nil nil black)*",
		"(1 nil 2 nil black)*(2 nil nil 1 red)",
		"(1 nil nil 2 red)(2 1 3 nil black)*(3 nil nil 2 red)",
		"(1 nil nil 2 black)(2 1 3 nil black)*(3 nil 4 2 black)(4 nil nil 3 red)",
		"(1 nil nil 2 black)(2 1 4 nil black)*(3 nil nil 4 red)(4 3 5 2 black)(5 nil nil 4 red)"}

	for i, td := range testData[:len(expected)] {
		tree := New()
		tree.Insert(td...)
		str := tree.String()
		if expected[i] != str {
			t.Errorf("\nexpected %v\nactual   %v", expected[i], str)
		}
	}

	var reverseExpected = []string{
		"(1 nil nil 2 red)(2 1 nil nil black)*",
		"(1 nil nil 2 red)(2 1 3 nil black)*(3 nil nil 2 red)",
		"(1 nil nil 2 red)(2 1 nil 3 black)(3 2 4 nil black)*(4 nil nil 3 black)"}

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
		sort.Sort(Int64Slice(td))
		if !reflect.DeepEqual(td, actual) {
			t.Error("Iterate failed on ", td)
		}
	}
}

func TestRangeReverse(t *testing.T) {
	for _, td := range reverseTestData {
		tree := New()
		tree.Insert(td...)
		if tree.Size() != len(td) {
			t.Error("Insert failed on input with ", len(td), " elements")
		}

		actual := make([]int64, 0, len(td))
		for n := range tree.Iterate() {
			actual = append(actual, n)
		}

		sort.Sort(sort.Reverse(Int64Slice(actual)))
		if !reflect.DeepEqual(td, actual) {
			t.Error("Expected: ", td, "Actual: ", actual)
		}
	}
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

func TestFail(t *testing.T) {
	values := []int64{1,2,3,4,5,6,7,8,9}
	expected := []int64{1,2,3,4,5,6,7,8,9}
	tree := New()
	tree.Insert(values...)
	slice := tree.Slice()
	t.Log(slice)
	if !reflect.DeepEqual(slice, expected) {
		t.Error("Expected ", expected, " got ", slice)
	}
}

func TestFind(t *testing.T) {
	values := []int64{24, 53, 70, 12, 96, 67, 61, 88, 28, 37, 16}
	tree := New()
	tree.Insert(values...)

    if tree.Find(24) == false {
        t.Error("couldn't find 24")
    }
    if tree.Find(16) == false {
        t.Error("couldn't find 16")
    }

    if tree.Find(42) == true {
        t.Error("found a number that didn't exist 16")
    }
}

func scrambleTesting(t *testing.T, length int) {
	expected := make([]int64, length)
	tree := New()
	for i := int64(0); i != int64(length); i++ {
		expected[i] = i
	}

	scrambled := scramble(0, expected)
	tree.Insert(scrambled...)
	slice := tree.Slice()
	if !reflect.DeepEqual(slice, expected) {
		t.Error("\nscrambled(", length, "): expected ", expected, "\ngot ", slice)
	}

	clone := tree.Clone()
	if !reflect.DeepEqual(slice, clone.Slice()) {
	    t.Error("\nscrambled(", length, "): Slice of the cloned tree is not equal. \nslice: ", slice, "\nclone:", clone.Slice())
	}
}

func TestScramble(t *testing.T) {
	scrambleTesting(t, 10)
	scrambleTesting(t, 100)
	scrambleTesting(t, 1000)
	scrambleTesting(t, 100000)
	scrambleTesting(t, 10000000)
}

func BenchmarkInsert1k(b *testing.B) {
	tree := New()
	for i := 0; i != b.N; i++ {
		tree.Insert(int64(i))
	}
}
