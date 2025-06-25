package ifacecommand

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kmio11/codegen/generator"
	"github.com/kmio11/codegen/generator/model"
	"github.com/kmio11/codegen/generator/parser"
)

// Command implements the interface generation command
type Command struct {
	fs              *flag.FlagSet
	flagPkg         *string
	flagType        *string
	flagOut         *string
	flagOutPkg      *string
	flagSelfPkgPath *string
	flagName        *string
}

// New creates a new interface command
func New() *Command {
	c := &Command{}
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.flagPkg = c.fs.String("pkg", ".", "The package containing the struct to generate interface from.")
	c.flagType = c.fs.String("type", "", "The name of the struct type to generate interface from.")
	c.flagOut = c.fs.String("out", "", "Output file; defaults to stdout.")
	c.flagOutPkg = c.fs.String("outpkg", "", "Output package name; defaults to the same package specified by -pkg")
	c.flagSelfPkgPath = c.fs.String("selfpkg", "", "The full package import path of the output package.")
	c.flagName = c.fs.String("name", "", "The name of the generated interface; defaults to <StructName>Interface")

	return c
}

// Name returns the command name
func (c Command) Name() string {
	return "interface"
}

// Description returns the command description
func (c Command) Description() string {
	return "generate interface from struct"
}

// Usage prints usage information
func (c Command) Usage(cmd string) {
	fmt.Printf(`Usage:
	%s %s [flags]

Generate Go interface from struct methods.

Examples:
	# Generate interface from struct in current package
	%s %s -pkg . -type UserService -out user_interface.go

	# Generate interface with custom name
	%s %s -pkg ./service -type UserService -name UserServiceInterface

	# Generate to different package
	%s %s -pkg ./internal -type Handler -outpkg contracts -out ./contracts/handler.go

Flags:
`, cmd, c.Name(), cmd, c.Name(), cmd, c.Name(), cmd, c.Name())
	c.fs.PrintDefaults()
}

// Parse parses command line arguments
func (c *Command) Parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return err
	}

	// Validate required flags
	if *c.flagType == "" {
		return fmt.Errorf("-type flag is required")
	}

	return nil
}

// Execute runs the interface generation command
func (c *Command) Execute() int {
	// Parse the package
	targetPkg, targetIntf, err := c.parse(*c.flagType, *c.flagPkg)
	if err != nil {
		log.Println(err)
		return 1
	}

	// Override interface name if specified
	if *c.flagName != "" {
		// Create new interface with custom name
		methods := targetIntf.Methods()
		pkgInfo := targetIntf.Type().Pkg()
		targetIntf = model.NewInterface(*c.flagName, pkgInfo, methods)
	}

	// Create output file
	file := c.createInterfaceFile(targetPkg, targetIntf, *c.flagOut, *c.flagOutPkg, *c.flagSelfPkgPath)

	// Generate code
	g := &generator.Generator{}
	src := g.
		PrintHeader(c.Name()).
		Printf("// Interface generated from %s.%s", targetPkg.Path, targetIntf.Name()).
		NewLine().
		Printf("%s", file.PrintCode()).
		Format()

	// Output
	if file.Path() == "" {
		fmt.Println(string(src))
	} else {
		err := os.WriteFile(file.Path(), src, 0644)
		if err != nil {
			log.Printf("writing output: %s\n", err)
			return 1
		}
		fmt.Printf("File created successfully : %s\n", file.Path())
	}

	return 0
}

// parse parses the package and extracts the target struct as interface
func (c *Command) parse(typ string, pkg string) (*model.Package, *model.Interface, error) {
	// Create parser
	p := parser.NewParser(
		parser.OptLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)),
		parser.OptParseTarget([]string{typ}),
	)

	// Load package
	err := p.LoadPackage(pkg)
	if err != nil {
		return nil, nil, err
	}

	// Parse the package
	targetPkg, err := p.Parse()
	if err != nil {
		return nil, nil, err
	}

	// Check if interface was generated
	if len(targetPkg.Interfaces) == 0 {
		return nil, nil, fmt.Errorf("no interface generated from struct %s", typ)
	}

	targetIntf := targetPkg.Interfaces[0]

	return targetPkg, targetIntf, nil
}

// createInterfaceFile creates a file containing the generated interface
func (c *Command) createInterfaceFile(targetPkg *model.Package, targetIntf *model.Interface, outFile, outPkgName, selfPkgPath string) *model.File {
	// Determine output package
	outPkgPath := selfPkgPath
	if outPkgName != "" && outPkgPath == "" {
		outPkgPath = outPkgName
	}
	if outPkgName == "" {
		outPkgName = targetPkg.Name
		outPkgPath = targetPkg.Path
	}

	// Create file
	file := model.NewFile(outFile, outPkgName, outPkgPath, targetPkg.CopyDependencies())
	file.DependenciesTidy()

	// Add interface to file
	file.AddInterface(targetIntf)

	return file
}
