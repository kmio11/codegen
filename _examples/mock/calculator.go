package mock

// Calculator is a simple interface for basic math operations
//
//go:generate go run ../.. mock -pkg . -type Calculator -out calculator_mock_gen.go
type Calculator interface {
	Add(a, b int) int
	Subtract(a, b int) int
	Multiply(a, b int) int
	Divide(a, b int) (int, error)
}

// Storage is a simple generic interface for key-value operations
//
//go:generate go run ../.. mock -pkg . -type Storage -out storage_mock_gen.go
type Storage[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
	List() []K
}