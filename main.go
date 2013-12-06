package main

import (
	"fmt"
	"os"
	"time"
)

type DepNode2I interface {
	Update()
	markAsNeedToUpdate()
	addSink(*DepNode2)
	getNeedToUpdate() bool
	markAsNotNeedToUpdate()
	getSources() []DepNode2I
}

type DepNode2 struct {
	needToUpdate bool
	sources      []DepNode2I
	sinks        []*DepNode2
}

func (this *DepNode2) InitDepNode2(sources ...DepNode2I) {
	this.needToUpdate = true
	this.sources = sources
	for _, source := range sources {
		source.addSink(this)
	}
}

func (this *DepNode2) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
}

func ForceUpdate(this DepNode2I) {
	this.markAsNeedToUpdate()
	MakeUpdated(this)
}

func MakeUpdated(this DepNode2I) {
	if !this.getNeedToUpdate() {
		return
	}
	for _, source := range this.getSources() {
		MakeUpdated(source)
	}
	this.Update()
	this.markAsNotNeedToUpdate()
}

func (this *DepNode2) markAsNeedToUpdate() {
	this.needToUpdate = true
	for _, sink := range this.sinks {
		sink.markAsNeedToUpdate()
	}
}

func (this *DepNode2) getNeedToUpdate() bool {
	return this.needToUpdate
}

func (this *DepNode2) markAsNotNeedToUpdate() {
	this.needToUpdate = false
}

func (this *DepNode2) getSources() []DepNode2I {
	return this.sources
}

// ---

type Node struct {
	Value int
	DepNode2
	name byte // Debug
}

func (n *Node) Update() {
	n.Value++
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

type NodeAdder Node

func (n *NodeAdder) Update() {
	n.Value = 0
	for _, source := range n.sources {
		n.Value += source.(*Node).Value
	}
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

type NodeMultiplier Node

func (n *NodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.sources {
		n.Value *= source.(*NodeAdder).Value
	}
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

var nodeA = &Node{name: 'A'}
var nodeB = &Node{name: 'B'}
var nodeT = &Node{name: 'T'}
var nodeX = &NodeAdder{name: 'X'}
var nodeY = &NodeAdder{name: 'Y'}
var nodeZ = &NodeMultiplier{name: 'Z'}
var Zlive = false

func main() {
	nodeX.InitDepNode2(nodeA, nodeB)
	nodeY.InitDepNode2(nodeB, nodeT)
	nodeZ.InitDepNode2(nodeX, nodeY)

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
				ForceUpdate(nodeA)
			case 'b':
				ForceUpdate(nodeB)
			case 'z':
				Zlive = !Zlive
				fmt.Println("Zlive changed to", Zlive) // Debug
			}
		case <-tick:
			ForceUpdate(nodeT)
		default:
		}

		if Zlive {
			MakeUpdated(nodeZ)
		}

		time.Sleep(5 * time.Millisecond)
	}
}
