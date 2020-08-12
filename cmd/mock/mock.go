package main

import (
	"codegen/pkg/generator"
	"codegen/pkg/generator/model"
	"codegen/pkg/generator/parser"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	cmdName = "genmock"
)

var (
	flagPkg         = flag.String("pkg", "", "The package containing interfaces to be mocked.")
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
	g.AddContents(file)

	g.PrintHeader(cmdName)
	g.Printf("// Mock for %s.%s", targetPkg.Path, targetIntf.Name)
	g.NewLine()
	g.PrintContents()
	src := g.Format()

	// output
	if len(*flagOut) == 0 {
		fmt.Println(string(src))
	} else {
		err := ioutil.WriteFile(file.Path, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
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
	file.ImportsTidy()

	// create mock impl
	mockImpl := mockImpl(targetPkg, targetIntf, outPkg)
	file.AddStruct(mockImpl)
	file.ImportsTidy()

	return file
}

func mockImpl(targetPkg *model.Package, targetIntf *model.Interface, outPkg *model.PkgInfo) *model.Struct {
	//mock struct
	mockName := "Mock" + targetIntf.Name
	mockImpl := model.NewStruct(mockName)

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
				targetIntf.Name,
			),
			"",
		),
	)

	for _, intfMethod := range targetIntf.Methods { //TODO: need to sort?
		// FakeFunction
		fakeFuncName := "Fake" + intfMethod.Name
		mockImpl.AddField(
			model.NewField(
				fakeFuncName,
				intfMethod.Type,
				"",
			),
		)

		// method
		methodRcvName := "m"
		methodRcv := *model.NewParameter(
			methodRcvName,
			model.NewTypeNamed(
				outPkg,
				mockImpl.Name,
			),
		)

		// method body
		/*
			return FakeXxx(x, x, x)
		*/
		var bodyCallFmt string
		if len(intfMethod.Type.Results) != 0 {
			bodyCallFmt += "return "
		}
		bodyCallFmt += methodRcvName + "." + fakeFuncName + intfMethod.Type.PrintCallArgsFmt()

		bodyCallArgs := []interface{}{}
		for _, a := range intfMethod.Type.Params {
			bodyCallArgs = append(bodyCallArgs, a.Name)
		}
		if intfMethod.Type.Variadic != nil {
			bodyCallArgs = append(bodyCallArgs, intfMethod.Type.Variadic.Name+"...")
		}
		methodBody := fmt.Sprintf(bodyCallFmt, bodyCallArgs...)

		// add method
		mockImpl.AddMethod(
			model.NewMethod(
				methodRcv,
				intfMethod.Name,
				intfMethod.Type,
				methodBody,
			),
		)
	}
	return mockImpl
}
