package main

import (
	"fmt"
	"os"
	"time"
)

type DepNode2I interface {
	MarkAsNeedToUpdate()
	Update()
	AddSink(*DepNode2)
	MakeUpdated()
}

type DepNode2 struct {
	NeedToUpdate bool
	Self         DepNode2I
	Sources      []DepNode2I
	Sinks        []*DepNode2
}

func (this *DepNode2) InitDepNode2(self DepNode2I, sources []DepNode2I) {
	this.NeedToUpdate = true
	this.Self = self
	this.Sources = sources
	for _, source := range sources {
		source.AddSink(this)
	}
}

func (this *DepNode2) AddSink(sink *DepNode2) {
	this.Sinks = append(this.Sinks, sink)
}

func (this *DepNode2) MakeUpdated() {
	if !this.NeedToUpdate || this.Self == nil {
		return
	}
	for _, source := range this.Sources {
		source.MakeUpdated()
	}
	this.Self.Update()
	this.NeedToUpdate = false
}

func (this *DepNode2) MarkAsNeedToUpdate() {
	this.NeedToUpdate = true
	for _, sink := range this.Sinks {
		sink.MarkAsNeedToUpdate()
	}
}

// ---

type Node struct {
	Value int
	DepNode2
	name byte // Debug
}

func (n *Node) Update() {
	fmt.Println("Updated", string(n.name), n.Value) // Debug
}

type NodeAdder Node

func (n *NodeAdder) Update() {
	n.Value = 0
	for _, source := range n.Sources {
		n.Value += source.(*Node).Value
	}
	fmt.Println("Updated", string(n.name), n.Value) // Debug
}

type NodeMultiplier Node

func (n *NodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.Sources {
		n.Value *= source.(*NodeAdder).Value
	}
	fmt.Println("Updated", string(n.name), n.Value) // Debug
}

var nodeA = &Node{name: 'A'}
var nodeB = &Node{name: 'B'}
var nodeT = &Node{name: 'T'}
var nodeX = &NodeAdder{name: 'X'}
var nodeY = &NodeAdder{name: 'Y'}
var nodeZ = &NodeMultiplier{name: 'Z'}
var Zlive = false

func main() {
	nodeX.InitDepNode2(nodeX, []DepNode2I{nodeA, nodeB})
	nodeY.InitDepNode2(nodeY, []DepNode2I{nodeB, nodeT})
	nodeZ.InitDepNode2(nodeZ, []DepNode2I{nodeX, nodeY})

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
				nodeA.MarkAsNeedToUpdate()
				nodeA.Update() // Debug
			case 'b':
				nodeB.Value++
				nodeB.MarkAsNeedToUpdate()
				nodeB.Update() // Debug
			case 'z':
				Zlive = !Zlive
				fmt.Println("Zlive changed to", Zlive) // Debug
			}
		case <-tick:
			nodeT.Value++
			nodeT.MarkAsNeedToUpdate()
			nodeT.Update() // Debug
		default:
		}

		if Zlive {
			nodeZ.MakeUpdated()
		}

		time.Sleep(5 * time.Millisecond)
	}
}
