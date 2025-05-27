package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the Go source file (replace "example.go" with your file)
	f, err := parser.ParseFile(fset, "../../examples/signup/business/index.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	fmt.Println("Global Variables:")
	//Global variables are declared at the package level as var declarations, represented as *ast.GenDecl with Tok == token.VAR. Their names are in ValueSpec.Names.
	for _, decl := range f.Decls {
		// Global variables
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						fmt.Println(" -", name.Name)
					}
				}
			}
		}
	}

	fmt.Println("\nFunction Names:")
	// Inspect the AST nodes
	ast.Inspect(f, func(n ast.Node) bool {
		// For example, print all function names
		if fn, ok := n.(*ast.FuncDecl); ok {
			fmt.Println(" -", fn.Name.Name)
		}
		return true
	})

	fmt.Println("\nFunction Arguments:")
	//Function arguments are found in FuncDecl.Type.Params.List. Each parameter may have multiple names (e.g., func f(a, b int)), so iterate over all names.
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			fmt.Printf(" %s args:\n", funcDecl.Name.Name)
			if funcDecl.Type.Params != nil {
				for _, param := range funcDecl.Type.Params.List {
					for _, name := range param.Names {
						fmt.Println("  -", name.Name)
					}
				}
			}
		}
	}
}
