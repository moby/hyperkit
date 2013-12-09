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

	Debug()
}

type DepNode2ManualI interface {
	DepNode2I
	manual() // Noop, just to separate it from automatic DepNode2I
}

// Updates dependencies and itself, only if its dependencies have changed.
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

// Updates dependencies and itself, regardless.
/*func ForceUpdated(this DepNode2I) {
	this.markAllAsNeedToUpdate()
	MakeUpdated(this)
}*/

// Updates only itself, regardless (skipping Update()).
func ExternallyUpdated(this DepNode2ManualI) {
	this.markAllAsNeedToUpdate()
	//this.markAsNotNeedToUpdate()
}

// ---

type DepNode2 struct {
	updated bool
	sources []DepNode2I
	sinks   []*DepNode2
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.sources
}

func (this *DepNode2) AddSources(sources ...DepNode2I) {
	this.updated = false
	this.sources = append(this.sources, sources...)
	for _, source := range sources {
		source.addSink(this)
	}
}

func (this *DepNode2) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
}

func (this *DepNode2) getNeedToUpdate() bool {
	return !this.updated
}

func (this *DepNode2) markAllAsNeedToUpdate() {
	this.updated = false
	for _, sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}

func (this *DepNode2) markAsNotNeedToUpdate() {
	this.updated = true
}

func (this *DepNode2) Debug() {
	fmt.Printf("%#v\n", this)
}

// ---

type DepNode2Manual struct {
	sinks []*DepNode2
}

func (this *DepNode2Manual) Update()                 { panic("") }
func (this *DepNode2Manual) GetSources() []DepNode2I { panic("") }
func (this *DepNode2Manual) addSink(sink *DepNode2) {
	this.sinks = append(this.sinks, sink)
}
func (this *DepNode2Manual) getNeedToUpdate() bool { return false }
func (this *DepNode2Manual) markAllAsNeedToUpdate() {
	for _, sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}
func (this *DepNode2Manual) markAsNotNeedToUpdate() { panic("") }
func (this *DepNode2Manual) manual()                { panic("") }
func (this *DepNode2Manual) Debug()                 { fmt.Printf("%#v\n", this) }

// ---

type DepNode2Func struct {
	UpdaterFunc func()
	DepNode2
}

func (this *DepNode2Func) Update() {
	this.UpdaterFunc()
}

// ---

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
	for _, source := range n.sources {
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
	for _, source := range n.sources {
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
