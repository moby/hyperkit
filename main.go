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

type DepNode2Blank struct {
	DepNode2
}

func (*DepNode2Blank) Update() {}

type DepNode2 struct {
	needToUpdate bool
	sources      []DepNode2I
	sinks        []*DepNode2
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.sources
}

func (this *DepNode2) AddSources(sources ...DepNode2I) {
	this.needToUpdate = true
	this.sources = sources
	for _, source := range sources {
		source.addSink(this)
	}
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

func ExternallyUpdated(this DepNode2I) {
	this.markAllAsNeedToUpdate()
	//MakeUpdated(this)
	this.markAsNotNeedToUpdate()
}

func (this *DepNode2) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
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

// ---

type node struct {
	Value int
	DepNode2
	name byte // Debug
}

// Debug
func (n *node) String() string {
	return fmt.Sprintf("%s -> %d", string(n.name), n.Value)
}

func (n *node) Update() {
	fmt.Println("Auto Updated", n) // Debug
}

type nodeAdder struct {
	node
}

func (n *nodeAdder) Update() {
	n.Value = 0
	for _, source := range n.sources {
		n.Value += source.(*node).Value
	}
	fmt.Println("Auto Updated", n) // Debug
}

type nodeMultiplier struct {
	node
}

func (n *nodeMultiplier) Update() {
	n.Value = 1
	for _, source := range n.sources {
		n.Value *= source.(*nodeAdder).Value
	}
	fmt.Println("Auto Updated", n) // Debug
}

var nodeA = &node{name: 'A'}
var nodeB = &node{name: 'B'}
var nodeT = &node{name: 'T'}
var nodeX = &nodeAdder{node{name: 'X'}}
var nodeY = &nodeAdder{node{name: 'Y'}}
var nodeZ = &nodeMultiplier{node{name: 'Z'}}
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
				ExternallyUpdated(nodeA)
				fmt.Println("User Updated", nodeA) // Debug
			case 'b':
				nodeB.Value++
				ExternallyUpdated(nodeB)
				fmt.Println("User Updated", nodeB) // Debug
			case 'z':
				zLive = !zLive
				fmt.Println("Zlive changed to", zLive) // Debug
			}
		case <-tick:
			nodeT.Value++
			ExternallyUpdated(nodeT)
			fmt.Println("Timer Updated", nodeT) // Debug
		default:
		}

		if zLive {
			MakeUpdated(nodeZ)
		}

		time.Sleep(5 * time.Millisecond)
	}
}
