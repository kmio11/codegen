package sample

// Basic generic interface with any constraint
type Repository[T any] interface {
	Get(id string) (T, error)
	Save(item T) error
	List() ([]T, error)
}

// Multiple type parameters
type Converter[From, To any] interface {
	Convert(from From) (To, error)
}

// Comparable constraint
type Comparer[T comparable] interface {
	Compare(a, b T) int
	Equal(a, b T) bool
}

// Union type constraint
type Number[T int | float64] interface {
	Add(a, b T) T
	Multiply(a, b T) T
}

// Complex interface constraint
type Stringable interface {
	String() string
}

type Formatter[T Stringable] interface {
	Format(item T) string
	FormatList(items []T) []string
}
