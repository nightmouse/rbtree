package rbtree

import (
    "testing"
    "fmt"
)

func TestNew(t *testing.T) {
    tree := New()
    if tree.String() != "" { 
        t.Error("empty tree somehow made a non-empty string")
    }
}

func TestInsert(t *testing.T) { 
    tree := New()
    tree.Insert(5,4,7,3,2,6,8,9)
    fmt.Println(tree)  
}
