package controller

import (
	batchv1 "github.com/sujeshthekkepatt/absurd-ci/api/v1"
)

type Node struct {
	Next *Node
	Val  batchv1.AStep
}

type List struct {
	Head   *Node
	Length int
}

func NewList() *List {

	return &List{}

}

func (l *List) Add(val batchv1.AStep) *List {

	node := &Node{
		Next: nil,
		Val:  val,
	}

	if l.Head == nil {

		l.Head = node
		l.Length = l.Length + 1
		return l
	} else {

		last := l.Head

		for last.Next != nil {
			last = last.Next
		}
		last.Next = node
		l.Length = l.Length + 1
		return l
	}
}

func (l *List) Get(val string) *Node {

	next := l.Head

	for next.Next != nil {
		if next.Val.Name == val {
			break
		} else {
			next = next.Next
		}

	}
	if next.Val.Name == val {
		return next
	} else {
		return nil
	}
}
