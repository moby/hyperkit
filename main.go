package gist7729255

type String interface {
	Get() string
}

type Strings interface {
	Get() []string
}

type StringsFunc func() []string

func (this StringsFunc) Get() []string {
	return this()
}
