package gist7802150

import (
	"fmt"
	"strings"
)

type DepNode2I interface {
	Update()

	GetSources() []DepNode2I

	addSink(*DepNode2)
	//removeSource(DepNode2I)
	getNeedToUpdate() bool
	markAllAsNeedToUpdate()
	markAsNotNeedToUpdate()
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
	sinks   map[*DepNode2]bool
}

func (this *DepNode2) GetSources() []DepNode2I {
	return this.sources
}

func (this *DepNode2) AddSources(sources ...DepNode2I) {
	//fmt.Println("AddSources", sources)
	this.updated = false
	this.sources = append(this.sources, sources...)
	for _, source := range sources {
		source.addSink(this)
	}
}

func (this *DepNode2) Close() error {
	/*for _, source := range this.sources {
		source.removeSink(this)
	}*/
	return nil
}

func (this *DepNode2) addSink(sink *DepNode2) {
	if this.sinks == nil {
		this.sinks = make(map[*DepNode2]bool)
	}
	this.sinks[sink] = true
}

func (this *DepNode2) removeSource(source DepNode2I) {
	for i, s := range this.sources {
		if s == source {
			this.sources = append(this.sources[:i], this.sources[i+1:]...)
			return
		}
	}
}

func (this *DepNode2) getNeedToUpdate() bool {
	return !this.updated
}

func (this *DepNode2) markAllAsNeedToUpdate() {
	this.updated = false
	for sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}

func (this *DepNode2) markAsNotNeedToUpdate() {
	this.updated = true
}

// ---

type DepNode2Manual struct {
	sinks map[*DepNode2]bool
}

func (this *DepNode2Manual) Update()                 { panic("") }
func (this *DepNode2Manual) GetSources() []DepNode2I { panic("") }
func (this *DepNode2Manual) Close() error {
	for sink := range this.sinks {
		sink.removeSource(this)
	}
	this.sinks = nil
	return nil
}
func (this *DepNode2Manual) addSink(sink *DepNode2) {
	if this.sinks == nil {
		this.sinks = make(map[*DepNode2]bool)
	}
	this.sinks[sink] = true
}

//func (this *DepNode2Manual) removeSource(source DepNode2I) { panic("") }
func (this *DepNode2Manual) getNeedToUpdate() bool { return false }
func (this *DepNode2Manual) markAllAsNeedToUpdate() {
	for sink := range this.sinks {
		// TODO: See if this can be optimized away...
		sink.markAllAsNeedToUpdate()
	}
}
func (this *DepNode2Manual) markAsNotNeedToUpdate() { panic("") }
func (this *DepNode2Manual) manual()                { panic("") }

// Given there are two distinct DepNode2Manual structs, each having a pointer,
// merge takes other and merges it (along with its current sinks) into this.
// Afterwards, both pointers point to a single unified DepNode2Manual struct.
func (this *DepNode2Manual) merge(other **DepNode2Manual) {
	if this.sinks == nil {
		this.sinks = make(map[*DepNode2]bool)
	}
	for sink := range (*other).sinks {
		this.sinks[sink] = true
	}

	*other = this
}

// ---

type DepNode2Func struct {
	UpdateFunc func(DepNode2I)
	DepNode2
}

func (this *DepNode2Func) Update() {
	this.UpdateFunc(this)
}

// =====

type ViewGroupI interface {
	SetSelf(string)

	AddAndSetViewGroup(ViewGroupI, string)
	RemoveView(ViewGroupI)

	GetUri() FileUri
	GetAllUris() []FileUri
	GetUriForProtocol(protocol string) (uri FileUri, ok bool)
	ContainsUri(FileUri) bool

	getViewGroup() *ViewGroup

	Debug()

	DepNode2ManualI
}

type FileUri string

type ViewGroup struct {
	all *map[ViewGroupI]bool
	uri FileUri

	*DepNode2Manual
}

func (this *ViewGroup) getViewGroup() *ViewGroup {
	return this
}

// InitViewGroup must be called after creating a new ViewGroupI,
// before any other ViewGroup method or ViewGroupI func.
func (this *ViewGroup) InitViewGroup(self ViewGroupI, uri FileUri) {
	this.all = &map[ViewGroupI]bool{self: true}
	this.uri = uri
	this.DepNode2Manual = &DepNode2Manual{}
}

// AddAndSetViewGroup adds another ViewGroupI and sets it to thisCurrent value, the current state of this ViewGroup.
func (this *ViewGroup) AddAndSetViewGroup(other ViewGroupI, thisCurrent string) {
	// Set other ViewGroup to thisCurrent
	for v := range *other.getViewGroup().all {
		v.SetSelf(thisCurrent)
	}
	ExternallyUpdated(other.getViewGroup().DepNode2Manual) // Notify whatever depended on the other ViewGroupI that it's been updated

	(*this.all)[other] = true
	other.getViewGroup().all = this.all
	this.DepNode2Manual.merge(&other.getViewGroup().DepNode2Manual)
}

// RemoveView removes a single view from the ViewGroup.
func (this *ViewGroup) RemoveView(other ViewGroupI) {
	delete(*this.all, other)
	//this.DepNode2Manual.Close() // <<< THIS is the key to making it work. But it's hack, need a proper solution.
	other.getViewGroup().InitViewGroup(other, other.GetUri())
}

func (this *ViewGroup) Debug() {
	fmt.Println(*this.all)
}

func (this *ViewGroup) GetUri() FileUri {
	return this.uri
}
func (this *ViewGroup) GetAllUris() (uris []FileUri) {
	for v := range *this.all {
		uris = append(uris, v.GetUri())
	}
	return uris
}
func (this *ViewGroup) GetUriForProtocol(protocol string) (uri FileUri, ok bool) {
	for v := range *this.all {
		if strings.HasPrefix(string(v.GetUri()), protocol) {
			return v.GetUri(), true
		}
	}
	return "", false
}
func (this *ViewGroup) ContainsUri(uri FileUri) bool {
	for v := range *this.all {
		if uri == v.GetUri() {
			return true
		}
	}
	return false
}

func SetViewGroup(this ViewGroupI, s string) {
	for v := range *this.getViewGroup().all {
		v.SetSelf(s)
	}

	ExternallyUpdated(this)
}

func SetViewGroupOther(this ViewGroupI, s string) {
	for v := range *this.getViewGroup().all {
		if v != this {
			v.SetSelf(s)
		}
	}

	ExternallyUpdated(this)
}
