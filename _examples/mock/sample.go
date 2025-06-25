package sample

// SomeInterface is interface (to be mocked.)
//
//go:generate go run github.com/kmio11/codegen/cmd/genmock -type SomeInterface -out testing_mock_gen.go
type SomeInterface interface {
	Name() string
	Add(x int, y int) (int, error)
}

// SomeImpl is impliment of SampleInterface
type SomeImpl struct {
	name string
}

var NewSomeImplFunc = func(name string) SomeInterface {
	return &SomeImpl{
		name: name,
	}
}

func NewSomeImpl(name string) SomeInterface {
	return NewSomeImplFunc(name)
}

func (s *SomeImpl) Name() string {
	return s.name
}

func (s *SomeImpl) Add(x int, y int) (int, error) {
	return x + y, nil
}
