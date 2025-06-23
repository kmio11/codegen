package mock

import (
	"flag"
	"io/ioutil"
	"sort"

	"github.com/kmio11/codegen/generator"
	"github.com/kmio11/codegen/generator/model"
	"github.com/kmio11/codegen/generator/parser"

	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Command is command
type Command struct {
	fs              *flag.FlagSet
	flagPkg         *string
	flagType        *string
	flagOut         *string
	flagOutPkg      *string
	flagSelfPkgPath *string
}

func New() *Command {
	c := &Command{}
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.flagPkg = c.fs.String("pkg", ".", "The package containing interfaces to be mocked.")
	c.flagType = c.fs.String("type", "", "The name of the type to be mocked.")
	c.flagOut = c.fs.String("out", "", "Output file; defaults to stdout.")
	c.flagOutPkg = c.fs.String("outpkg", "", "Output package name; defaults to the same package specified by -pkg")
	c.flagSelfPkgPath = c.fs.String("selfpkg", "", "The full package import path of the output package.")

	return c
}

func (c Command) Name() string {
	return "mock"
}

func (c Command) Description() string {
	return "generate mock"
}

func (c Command) Usage(cmd string) {
	fmt.Printf(`Usage:
    %s %s -pkg <package> -type <type> [-out <out>] [-outpkg <outpkg> [-selfpkg <selfpkg>]]

`,
		cmd, c.Name(),
	)
	c.fs.PrintDefaults()
}

func (c Command) Parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return err
	}
	if len(*c.flagPkg)*len(*c.flagType) == 0 {
		return fmt.Errorf("")
	}
	if len(*c.flagOutPkg) == 0 && len(*c.flagSelfPkgPath) != 0 {
		return fmt.Errorf("")
	}

	return nil
}

func (c Command) Execute() int {
	// parse
	targetPkg, targetIntf, err := parse(*c.flagType, *c.flagPkg)
	if err != nil {
		log.Println(err)
		return 1
	}

	// create mock
	file := mockfile(targetPkg, targetIntf, *c.flagOut, *c.flagOutPkg, *c.flagSelfPkgPath)

	// generate
	g := &generator.Generator{}
	src := g.
		PrintHeader(c.Name()).
		Printf("// Mock for %s.%s", targetPkg.Path, targetIntf.Name()).
		NewLine().
		Printf(file.PrintCode()).
		Format()

	// output
	if file.Path() == "" {
		fmt.Println(string(src))
	} else {
		err := ioutil.WriteFile(file.Path(), src, 0644)
		if err != nil {
			log.Printf("writing output: %s\n", err)
			return 1
		}
		fmt.Printf("File created successfully : %s\n", file.Path())
	}

	return 0
}

func parse(typ string, pkg string) (*model.Package, *model.Interface, error) {
	// parse target package
	parser := parser.NewParser(
		parser.OptLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)),
		parser.OptParseTarget([]string{typ}),
	)
	patterns := pkg
	err := parser.LoadPackage(patterns)
	if err != nil {
		return nil, nil, err
	}
	targetPkg, err := parser.Parse()
	if err != nil {
		return nil, nil, err

	}
	targetIntf := targetPkg.Interfaces[0]

	return targetPkg, targetIntf, nil
}

func mockfile(targetPkg *model.Package, targetIntf *model.Interface, outFile, outPkgName, selfPkgPath string) *model.File {
	// output
	outPkgPath := selfPkgPath
	if outPkgName != "" && outPkgPath == "" {
		outPkgPath = outPkgName
	}
	if outPkgName == "" {
		outPkgName = targetPkg.Name
		outPkgPath = targetPkg.Path
	}
	outPkg := model.NewPkgInfo(outPkgName, outPkgPath, "")

	// create file which has mock.
	file := model.NewFile(outFile, outPkgName, outPkgPath, targetPkg.CopyDependencies())
	file.DependenciesTidy()

	// create mock impl
	mockImpl := mockImpl(targetPkg, targetIntf, outPkg)
	file.AddStruct(mockImpl)

	// create stub
	stubRoot, stubs := stub(targetPkg, targetIntf, outPkg, mockImpl)

	file.AddStruct(stubRoot)
	for _, stub := range stubs {
		file.AddStruct(stub)
	}

	file.DependenciesTidy()
	return file
}

const (
	mockRcvName = "m"
	stubRcvName = "s"
)

func getMockFieldName(intfMethodName string) string {
	return "Fake" + intfMethodName
}

func getMockArgsName(i int) string {
	return "a" + strconv.Itoa(i)
}

func getMockResultsName(i int) string {
	return ""
}

func getStubMethodName(intfMethodName string) string {
	return getMockFieldName(intfMethodName)
}

// fmtSignature returns *mode.TypeSignature
// param and results names replaced by no duplication names.
func fmtSignature(org *model.TypeSignature) *model.TypeSignature {
	methodParams := []*model.Parameter{}
	var n int
	for _, p := range org.Args() {
		methodParams = append(methodParams,
			model.NewParameter(getMockArgsName(n), p.Type()),
		)
		n++
	}
	var methodVariadic *model.Parameter
	if org.Variadic() != nil {
		methodVariadic = model.NewParameter(getMockArgsName(n), org.Variadic().Type())
	}
	methodResults := []*model.Parameter{}
	for i, r := range org.Results() {
		methodResults = append(methodResults,
			model.NewParameter(getMockResultsName(i), r.Type()),
		)
	}
	return model.NewTypeSignature(
		methodParams,
		methodVariadic,
		methodResults,
	)
}

func mockImpl(targetPkg *model.Package, targetIntf *model.Interface, outPkg *model.PkgInfo) *model.Struct {
	//mock struct
	mockName := "Mock" + targetIntf.Name()
	var mockImpl *model.Struct

	// Handle generic interfaces
	if targetIntf.IsGeneric() {
		mockImpl = model.NewGenericStruct(mockName, outPkg, targetIntf.TypeParams())
	} else {
		mockImpl = model.NewStruct(mockName, outPkg)
	}

	// add interface to mock
	var interfaceType *model.TypeNamed
	if targetIntf.IsGeneric() {
		// For embedded generic interfaces, use type parameters without constraints
		typeParamsWithoutConstraints := []*model.TypeParameter{}
		for i, param := range targetIntf.TypeParams() {
			typeParamsWithoutConstraints = append(typeParamsWithoutConstraints,
				model.NewTypeParameter(param.Name(), nil, i))
		}
		interfaceType = model.NewGenericTypeNamed(
			model.NewPkgInfo(
				targetPkg.Name,
				targetPkg.Path,
				"",
			),
			targetIntf.Name(),
			targetIntf.Type().Org(),
			typeParamsWithoutConstraints,
		)
	} else {
		interfaceType = model.NewTypeNamed(
			model.NewPkgInfo(
				targetPkg.Name,
				targetPkg.Path,
				"",
			),
			targetIntf.Name(),
			targetIntf.Type().Org(),
		)
	}

	mockImpl.AddField(
		model.NewField(
			"",
			interfaceType,
			"",
		),
	)

	for _, intfMethod := range targetIntf.Methods() {
		// Mock's Fields: FakeFunction
		fakeFuncName := getMockFieldName(intfMethod.Name())
		mockImpl.AddField(
			model.NewField(
				fakeFuncName,
				intfMethod.Type(),
				"",
			),
		)

		// Mock's methods
		// For method receivers on generic types, include type parameters without constraints
		var baseType *model.TypeNamed
		if mockImpl.IsGeneric() {
			// Create type parameters without constraints for method receivers
			typeParamsNoConstraints := []*model.TypeParameter{}
			for i, param := range mockImpl.TypeParams() {
				typeParamsNoConstraints = append(typeParamsNoConstraints,
					model.NewTypeParameter(param.Name(), nil, i))
			}
			baseType = model.NewGenericTypeNamed(
				outPkg,
				mockImpl.Name(),
				mockImpl.TypeStruct(),
				typeParamsNoConstraints,
			)
		} else {
			baseType = model.NewTypeNamed(
				outPkg,
				mockImpl.Name(),
				mockImpl.TypeStruct(),
			)
		}
		methodRcvType := model.NewPointer(baseType)

		methodRcv := model.NewParameter(
			mockRcvName,
			methodRcvType,
		)

		// method body
		/*
			return FakeXxx(a0, a1, a2)
		*/
		var bodyCallFmt string
		if len(intfMethod.Type().Results()) != 0 {
			bodyCallFmt += "return "
		}
		bodyCallFmt += mockRcvName + "." + fakeFuncName + intfMethod.Type().PrintCallArgsFmt()

		bodyCallArgs := []interface{}{}
		var n int
		for range intfMethod.Type().Args() {
			bodyCallArgs = append(bodyCallArgs, getMockArgsName(n))
			n++
		}
		if intfMethod.Type().Variadic() != nil {
			bodyCallArgs = append(bodyCallArgs, getMockArgsName(n)+"...")
		}
		methodBody := fmt.Sprintf(bodyCallFmt, bodyCallArgs...)

		// add method
		mockImpl.AddMethod(
			model.NewMethod(
				methodRcv,
				intfMethod.Name(),
				fmtSignature(intfMethod.Type()),
				methodBody,
			),
		)
	}
	return mockImpl
}

func stub(targetPkg *model.Package, targetIntf *model.Interface, outPkg *model.PkgInfo, mockImpl *model.Struct) (stubRoot *model.Struct, stubs []*model.Struct) {
	stubRootName := "Stub" + targetIntf.Name()

	// Handle generic interfaces for stub root
	if targetIntf.IsGeneric() {
		stubRoot = model.NewGenericStruct(stubRootName, outPkg, targetIntf.TypeParams())
	} else {
		stubRoot = model.NewStruct(stubRootName, outPkg)
	}

	// For method receivers on generic types, include type parameters without constraints
	var stubRootBaseType *model.TypeNamed
	if stubRoot.IsGeneric() {
		// Create type parameters without constraints for method receivers
		typeParamsNoConstraints := []*model.TypeParameter{}
		for i, param := range stubRoot.TypeParams() {
			typeParamsNoConstraints = append(typeParamsNoConstraints,
				model.NewTypeParameter(param.Name(), nil, i))
		}
		stubRootBaseType = model.NewGenericTypeNamed(outPkg, stubRootName, stubRoot.TypeStruct(), typeParamsNoConstraints)
	} else {
		stubRootBaseType = model.NewTypeNamed(outPkg, stubRootName, stubRoot.TypeStruct())
	}
	stubRootRcv := model.NewParameter(stubRcvName, model.NewPointer(stubRootBaseType))
	stubMethods := []*model.Method{}

	mockInitVals := map[string]string{} // for NewMockBody

	stubs = []*model.Struct{}
	for _, intfMethod := range targetIntf.Methods() {
		// stub for each intf's method.
		stubName := "Stub" + intfMethod.Name()
		var stub *model.Struct

		// Handle generic interfaces for individual stub structs
		if targetIntf.IsGeneric() {
			stub = model.NewGenericStruct(stubName, outPkg, targetIntf.TypeParams())
		} else {
			stub = model.NewStruct(stubName, outPkg)
		}
		for i, param := range intfMethod.Type().Results() {
			stub.AddField(
				model.NewField(
					"R"+strconv.Itoa(i),
					param.Type(),
					"",
				),
			)
		}
		stubs = append(stubs, stub)

		// stubRoot's method for each intf's method.
		stubFiealdName := intfMethod.Name()
		stubRoot.AddField(
			model.NewField(
				stubFiealdName,
				stub.Type(),
				"",
			),
		)

		stubMethodBody := "return "
		for j, r := range stub.Fields() {
			if j > 0 {
				stubMethodBody += ","
			}
			stubMethodBody += stubRootRcv.Name() + "." + stubFiealdName + "." + r.Name()
		}
		stubMethodName := getStubMethodName(intfMethod.Name())
		stubMethods = append(stubMethods,
			model.NewMethod(
				stubRootRcv,
				stubMethodName,
				fmtSignature(intfMethod.Type()),
				stubMethodBody,
			),
		)

		// for NewMock
		mockInitVals[getMockFieldName(intfMethod.Name())] = stubRootRcv.Name() + "." + stubMethodName
	}

	// NewMock method
	// sort
	keys := []string{}
	for k := range mockInitVals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var newMockBody string
	var returnType *model.TypeNamed

	if targetIntf.IsGeneric() {
		// For generic interfaces, include type parameters in mock creation (without constraints)
		mockTypeName := mockImpl.Name() + "["
		for i, param := range targetIntf.TypeParams() {
			if i > 0 {
				mockTypeName += ", "
			}
			mockTypeName += param.Name()
			// Don't include constraints in instantiation
		}
		mockTypeName += "]"
		newMockBody = "return &" + mockTypeName + "{"

		// Create generic interface return type
		returnType = model.NewGenericTypeNamed(
			model.NewPkgInfo(
				targetPkg.Name,
				targetPkg.Path,
				"",
			),
			targetIntf.Name(),
			targetIntf.Type().Org(),
			targetIntf.TypeParams(),
		)
	} else {
		newMockBody = "return &" + mockImpl.Name() + "{"
		returnType = targetIntf.Type()
	}

	for _, k := range keys {
		v := mockInitVals[k]
		newMockBody += k + ":" + v
		newMockBody += ","
	}
	newMockBody = strings.TrimRight(newMockBody, ",")
	newMockBody += "}"

	newMock := model.NewMethod(
		stubRootRcv,
		"NewMock",
		model.NewTypeSignature(nil, nil,
			[]*model.Parameter{
				model.NewParameter("", returnType),
			},
		),
		newMockBody,
	)

	// add method to stubRoot
	stubRoot.AddMethod(newMock)
	for _, m := range stubMethods {
		stubRoot.AddMethod(m)
	}

	return
}
