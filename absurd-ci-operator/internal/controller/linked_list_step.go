package controller

import batchv1 "github.com/sujeshthekkepatt/absurd-ci/api/v1"

type LinkedList struct {
}

type Node struct {
	Next  *Node
	Value batchv1.AStep
}
