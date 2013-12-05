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
	GetNeedToUpdate() bool
	MarkAsNotNeedToUpdate()
	GetSources() []DepNode2I
}

type DepNode2 struct {
	NeedToUpdate bool
	Sources      []DepNode2I
	Sinks        []*DepNode2
}

func (this *DepNode2) InitDepNode2(sources ...DepNode2I) {
	this.NeedToUpdate = true
	this.Sources = sources
	for _, source := range sources {
		source.AddSink(this)
	}
}

func (this *DepNode2) AddSink(sink *DepNode2) {
	this.Sinks = append(this.Sinks, sink)
}

func ForceUpdate(this DepNode2I) {
	this.MarkAsNeedToUpdate()
	MakeUpdated(this)
}

func MakeUpdated(this DepNode2I) {
	if !this.GetNeedToUpdate() {
		return
	}
	for _, source := range this.GetSources() {
		MakeUpdated(source)
	}
	this.Update()
	this.MarkAsNotNeedToUpdate()
}

func (this *DepNode2) MarkAsNeedToUpdate() {
	this.NeedToUpdate = true
	for _, sink := range this.Sinks {
		sink.MarkAsNeedToUpdate()
	}
}

func (this *DepNode2) GetNeedToUpdate() bool {
	return this.NeedToUpdate
}

func (this *DepNode2) MarkAsNotNeedToUpdate() {
	this.NeedToUpdate = false
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.Sources
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
	for _, source := range n.Sources {
		n.Value += source.(*Node).Value
	}
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

type NodeMultiplier Node

func (n *NodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.Sources {
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
