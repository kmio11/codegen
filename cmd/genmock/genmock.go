package main

import (
	"github.com/kmio11/codegen/generator"
	"github.com/kmio11/codegen/generator/model"
	"github.com/kmio11/codegen/generator/parser"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	cmdName = "genmock"
)

var (
	flagPkg         = flag.String("pkg", ".", "The package containing interfaces to be mocked.")
	flagType        = flag.String("type", "", "The name of the type to be mocked.")
	flagOut         = flag.String("out", "", "Output file; defaults to stdout.")
	flagOutPkg      = flag.String("outpkg", "", "Output package name; defaults to the same package specified by -pkg")
	flagSelfPkgPath = flag.String("selfpkg", "", "The full package import path of the output package.")
)

const usageTxt = `
usage: 
    mock -pkg <package> -type <type> [-out <out>] [-outpkg <outpkg> [-selfpkg <selfpkg>]]
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
		flag.PrintDefaults()
	}
}

func usage() {
	flag.Usage()
	os.Exit(2)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// parse args
	flag.Parse()
	if len(*flagPkg)*len(*flagType) == 0 {
		usage()
	}
	if len(*flagOutPkg) == 0 && len(*flagSelfPkgPath) != 0 {
		usage()
	}

	// parse
	targetPkg, targetIntf := parse()

	// create mock
	file := mockfile(targetPkg, targetIntf)

	// generate
	g := &generator.Generator{}
	src := g.
		PrintHeader(cmdName).
		Printf("// Mock for %s.%s", targetPkg.Path, targetIntf.Name()).
		NewLine().
		AddContents(file).PrintContents().
		Format()

	// output
	if len(*flagOut) == 0 {
		fmt.Println(string(src))
	} else {
		err := ioutil.WriteFile(file.Path(), src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s\n", err)
		}
		fmt.Printf("File created successfully : %s\n", file.Path())
	}

	os.Exit(0)
}

func parse() (*model.Package, *model.Interface) {
	// parse target package
	parser := parser.NewParser(
		parser.OptLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)),
		parser.OptParseTarget([]string{*flagType}),
	)
	patterns := *flagPkg
	err := parser.LoadPackage(patterns)
	if err != nil {
		log.Fatal(err)
	}
	targetPkg, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	targetIntf := targetPkg.Interfaces[0]

	return targetPkg, targetIntf
}

func mockfile(targetPkg *model.Package, targetIntf *model.Interface) *model.File {
	// output
	outFile := *flagOut
	outPkgName := *flagOutPkg
	outPkgPath := *flagSelfPkgPath
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
	mockImpl := model.NewStruct(mockName, outPkg)

	// add interface to mock
	mockImpl.AddField(
		model.NewField(
			"",
			model.NewTypeNamed(
				model.NewPkgInfo(
					targetPkg.Name,
					targetPkg.Path,
					"",
				),
				targetIntf.Name(),
				targetIntf.Type(),
			),
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
		methodRcv := model.NewParameter(
			mockRcvName,
			model.NewTypeNamed(
				outPkg,
				mockImpl.Name(),
				mockImpl.TypeStruct(),
			),
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
	stubRoot = model.NewStruct(stubRootName, outPkg)
	stubRootRcv := model.NewParameter(stubRcvName, stubRoot.Type())
	stubMethods := []*model.Method{}

	mockInitVals := map[string]string{} // for NewMockBody

	stubs = []*model.Struct{}
	for _, intfMethod := range targetIntf.Methods() {
		// stub for each intf's method.
		stubName := "Stub" + intfMethod.Name()
		stub := model.NewStruct(stubName, outPkg)
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
	newMockBody := "return &" + mockImpl.Name() + "{"
	for k, v := range mockInitVals {
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
				model.NewParameter("", targetIntf.Type()),
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
