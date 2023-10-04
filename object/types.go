package object

type MatchType int

const (
	ExactMatch MatchType = iota
	PrefixMatch
	SuffixMatch
	SubMatch
	RegExMatch
)

type DataTypes interface {
	Opaque | Refrence
}

type Object[T DataTypes] interface {
	Children() []string
	Child(any) Selector[T]
	Set(string, T)
	Get(string) T
	Map() map[string]any
	Match(string, MatchType) ([]Object[T], error)
}

type Selector[T DataTypes] interface {
	Name() string
	Exists() bool
	Set(string, T) error
	Get(string) (T, error)
	Add(Object[T]) error
	Object() (Object[T], error)
}

type Opaque []byte

type Refrence any
