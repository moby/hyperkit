package gist7802150

import ()

type DepNode2I interface {
	Update()
	GetSources() []DepNode2I

	addSink(*DepNode2)
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

// ---

type DepNode2Func struct {
	UpdaterFunc func()
	DepNode2
}

func (this *DepNode2Func) Update() {
	this.UpdaterFunc()
}

// =====

type ViewGroupI interface {
	SetSelf(string)

	getViewGroup() *ViewGroup

	DepNode2ManualI
}

type ViewGroup struct {
	all *map[ViewGroupI]bool

	*DepNode2Manual
}

func (this *ViewGroup) getViewGroup() *ViewGroup {
	return this
}

func (this *ViewGroup) InitViewGroup(self ViewGroupI) {
	this.all = &map[ViewGroupI]bool{self: true}
	this.DepNode2Manual = &DepNode2Manual{}
}

// TODO: Change to func JoinViewGroups(a, b ViewGroupI) to better match its symmetrical nature
func (this *ViewGroup) AddView(other ViewGroupI) {
	(*this.all)[other] = true
	other.getViewGroup().all = this.all
	other.getViewGroup().DepNode2Manual = this.DepNode2Manual
}

func (this *ViewGroup) RemoveView(other ViewGroupI) {
	delete(*this.all, other)
	other.getViewGroup().InitViewGroup(other)
}

func SetViewGroup(this ViewGroupI, s string) {
	if this.getViewGroup().all != nil {
		for v := range *this.getViewGroup().all {
			v.SetSelf(s)
		}
	}

	ExternallyUpdated(this)
}

func SetViewGroupOther(this ViewGroupI, s string) {
	if this.getViewGroup().all != nil {
		for v := range *this.getViewGroup().all {
			if v != this {
				v.SetSelf(s)
			}
		}
	}

	ExternallyUpdated(this)
}
