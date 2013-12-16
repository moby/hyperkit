package gist7729255

type String interface {
	Get() string
}

type StringFunc func() string

func (this StringFunc) Get() string {
	return this()
}

// ---

type Strings interface {
	Get() []string
}

type StringsFunc func() []string

func (this StringsFunc) Get() []string {
	return this()
}
