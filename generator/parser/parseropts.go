package parser

import (
	"io"
	"log"
)

// Opts is option of Parser.
type Opts func(*Parser)

// OptPackage is set package to paraser.
func OptPackage(pkg *Package) Opts {
	return func(p *Parser) {
		p.ParsedPkg = pkg
	}
}

// OptParseTarget set parse target to paraser.
func OptParseTarget(targets []string) Opts {
	return func(p *Parser) {
		p.Targets = targets
	}
}

// OptLogger set logger.
func OptLogger(logger *log.Logger) Opts {
	return func(p *Parser) {
		if logger == nil {
			p.log = log.New(io.Discard, "", 0)
			return
		}
		p.log = logger
	}
}

// OptStopLoadErr stop processing if the error that occurred when loading the package .
func OptStopLoadErr() Opts {
	return func(p *Parser) {
		p.stopLoadErr = true
	}
}
