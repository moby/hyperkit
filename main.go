package main

import (
	"fmt"
	"os"
	"time"
)

type DepNode2I interface {
	Update()
	GetSources() []DepNode2I
	addSink(*DepNode2)
	getNeedToUpdate() bool
	markAllAsNeedToUpdate()
	markAsNotNeedToUpdate()
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
	this.markAllAsNeedToUpdate()
	MakeUpdated(this)
}

func MakeUpdated(this DepNode2I) {
	if !this.getNeedToUpdate() {
		return
	}
	for _, source := range this.GetSources() {
		MakeUpdated(source)
	}
	this.Update()
	this.markAsNotNeedToUpdate()
}

func (this *DepNode2) markAllAsNeedToUpdate() {
	this.needToUpdate = true
	for _, sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}

func (this *DepNode2) getNeedToUpdate() bool {
	return this.needToUpdate
}

func (this *DepNode2) markAsNotNeedToUpdate() {
	this.needToUpdate = false
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.sources
}

// ---

type node struct {
	Value int
	DepNode2
	name byte // Debug
}

func (n *node) Update() {
	n.Value++
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

type nodeAdder node

func (n *nodeAdder) Update() {
	n.Value = 0
	for _, source := range n.sources {
		n.Value += source.(*node).Value
	}
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

type nodeMultiplier node

func (n *nodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.sources {
		n.Value *= source.(*nodeAdder).Value
	}
	fmt.Println("Updated", string(n.name), "to", n.Value) // Debug
}

var nodeA = &node{name: 'A'}
var nodeB = &node{name: 'B'}
var nodeT = &node{name: 'T'}
var nodeX = &nodeAdder{name: 'X'}
var nodeY = &nodeAdder{name: 'Y'}
var nodeZ = &nodeMultiplier{name: 'Z'}
var zLive = false

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
				zLive = !zLive
				fmt.Println("Zlive changed to", zLive) // Debug
			}
		case <-tick:
			ForceUpdate(nodeT)
		default:
		}

		if zLive {
			MakeUpdated(nodeZ)
		}

		time.Sleep(5 * time.Millisecond)
	}
}
