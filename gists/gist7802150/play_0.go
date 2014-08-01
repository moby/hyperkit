// +build ignore

package main

import . "gist.github.com/7802150.git"

import (
	"fmt"
	"os"
	"time"
)

type node struct {
	Value int
	DepNode2Manual
	name byte // Debug
}

// Debug
func (n *node) String() string {
	return fmt.Sprintf("%s -> %d", string(n.name), n.Value)
}

type nodeAdder struct {
	Value int
	DepNode2
	name byte // Debug
}

// Debug
func (n *nodeAdder) String() string {
	return fmt.Sprintf("%s -> %d", string(n.name), n.Value)
}

func (n *nodeAdder) Update() {
	n.Value = 0
	for _, source := range n.GetSources() {
		n.Value += source.(*node).Value
	}
	fmt.Println("Auto Updated", n) // Debug
}

type nodeMultiplier struct {
	Value int
	DepNode2
	name byte // Debug
}

// Debug
func (n *nodeMultiplier) String() string {
	return fmt.Sprintf("%s -> %d", string(n.name), n.Value)
}

func (n *nodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.GetSources() {
		n.Value *= source.(*nodeAdder).Value
	}
	fmt.Println("Auto Updated", n) // Debug
}

var nodeA = &node{name: 'A'}
var nodeB = &node{name: 'B'}
var nodeT = &node{name: 'T'}
var nodeX = &nodeAdder{name: 'X'}
var nodeY = &nodeAdder{name: 'Y'}
var nodeZ = &nodeMultiplier{name: 'Z'}
var zLive = false

func main() {
	nodeX.AddSources(nodeA, nodeB)
	nodeY.AddSources(nodeB, nodeT)
	nodeZ.AddSources(nodeX, nodeY)

	user := make(chan byte)
	go func() {
		b := make([]byte, 1024)
		for {
			os.Stdin.Read(b)
			user <- b[0]
		}
	}()

	tick := time.Tick(10 * time.Second)

	for {
		select {
		case c := <-user:
			switch c {
			case 'a':
				nodeA.Value++
				ExternallyUpdated(&nodeA.DepNode2Manual)
				fmt.Println("User Updated", nodeA) // Debug
			case 'b':
				nodeB.Value++
				ExternallyUpdated(&nodeB.DepNode2Manual)
				fmt.Println("User Updated", nodeB) // Debug
			case 'z':
				zLive = !zLive
				fmt.Println("Zlive changed to", zLive) // Debug
			}
		case <-tick:
			nodeT.Value++
			ExternallyUpdated(&nodeT.DepNode2Manual)
			fmt.Println("Timer Updated", nodeT) // Debug
		default:
		}

		if zLive {
			MakeUpdated(nodeZ)
		}

		time.Sleep(5 * time.Millisecond)
	}
}
