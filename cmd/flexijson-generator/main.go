package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var extraRx = regexp.MustCompile(`\bextrajson\b`)

type metadata struct {
	PackageName string
	Types       []typeMetadata
}

type typeMetadata struct {
	WrappingTypeName, TargetTypeName, ExtraFieldName string
}

var (
	cli = kingpin.New("flexijson-generator", "Generates a wrapping type for a struct to handle extra fields in JSON when unmarshalling/marshalling.")

	cliFlagPackageName = cli.Flag("package", "The target package name.").Short('p').Required().String()
	cliFlagSourcePath  = cli.Flag("input", "The input source code file.").Short('i').Required().ExistingFile()
	cliFlagOutputFile  = cli.Flag("output", "The output source code file.").Short('o').Required().String()
)

func main() {
	kingpin.MustParse(cli.Parse(os.Args[1:]))

	// Parse input go code
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, *cliFlagSourcePath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	output := new(bytes.Buffer)
	meta := metadata{
		Types: []typeMetadata{},
	}

	meta.PackageName = astFile.Name.Name

	// loop through parsed declarations
	for _, declaration := range astFile.Decls {
		switch typedDeclaration := declaration.(type) {
		case *ast.GenDecl: // import, constant, type or variable declaration
			// loop through specifications
			for _, spec := range typedDeclaration.Specs {
				switch typedSpec := spec.(type) {
				case *ast.TypeSpec: // type declaration
					typeMeta := typeMetadata{
						TargetTypeName:   typedSpec.Name.Name,
						WrappingTypeName: strings.ToUpper(typedSpec.Name.Name[0:1]) + typedSpec.Name.Name[1:],
					}

					// Find extra field name if any exists
					switch typedType := typedSpec.Type.(type) {
					case *ast.StructType:
						for _, field := range typedType.Fields.List {
							if field.Tag != nil {
								if extraRx.MatchString(field.Tag.Value) {
									// that's an extra field
									typeMeta.ExtraFieldName = field.Names[0].Name
									break
								}
							}
						}
						break
					}

					// Add type to metadata for template
					meta.Types = append(meta.Types, typeMeta)
					break
				}
			}
			break
		}
	}

	// generate output source code using template
	err = wrappingClassTemplate.Execute(output, meta)
	if err != nil {
		log.Fatal(err)
	}

	// format output source code
	formattedBytes, err := format.Source(output.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// write output source code to output file
	outputFile, err := os.Create(*cliFlagOutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(formattedBytes)
	if err != nil {
		log.Fatal(err)
	}
}
