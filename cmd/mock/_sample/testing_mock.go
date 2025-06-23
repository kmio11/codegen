package sample

// SetMock : This is not auto-generated.
func (s StubSomeInterface) SetMock() {
	NewSomeImplFunc = func(name string) SomeInterface {
		return s.NewMock()
	}
}
