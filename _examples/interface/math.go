package main

// Calculator is a simple calculator implementation
//
//go:generate go run ../.. interface -pkg . -type Calculator -out calculator_interface_gen.go
type Calculator struct{}

func (c *Calculator) Add(a, b int) int {
	return a + b
}

func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

func (c *Calculator) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, nil // simplified for example
	}
	return a / b, nil
}

// Logger is a simple logging implementation
//
//go:generate go run ../.. interface -pkg . -type Logger -out logger_interface_gen.go
type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(message string) {
	// implementation omitted for example
}

func (l *Logger) Error(message string) {
	// implementation omitted for example
}

func (l *Logger) Debug(message string) {
	// implementation omitted for example
}