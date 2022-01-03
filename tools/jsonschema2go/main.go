package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gebv/go-lsp/tools/jsonschema2go/internal"

	"github.com/pkg/errors"
	js "github.com/santhosh-tekuri/jsonschema/v5"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Want two args: workdir and pattern")
	}
	workDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	pattern := os.Args[2]

	// input args
	fmt.Println("Work Dir:", workDir)
	fmt.Println("Pattern:", pattern)

	files, _ := internal.ListFiles(workDir, pattern)
	fmt.Println("Num Files:", len(files))

	// process files
	for idx := range files {
		filePath := files[idx]
		fmt.Println("Process File:", filePath)
		if err := processFile(filePath); err != nil && err != internal.ErrCircularDep {
			log.Fatalf("Failed process file %q: %v\n", filePath, err)
		}
	}
}

func processFile(filePath string) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	jsc := js.NewCompiler()
	jsc.Draft = js.Draft7
	jsc.ExtractAnnotations = true
	// jsc.LoadURL = func(s string) (io.ReadCloser, error) {
	// 	fmt.Println("load", s)
	// 	return nil, nil
	// }

	err = jsc.AddResource("root", bytes.NewBuffer(fileBytes))
	if err != nil {
		return errors.Wrap(err, "add files to manage of resources")
	}
	s, err := jsc.Compile("root")
	if err != nil {
		return errors.Wrap(err, "compile")
	}

	// sid, err := internal.SchemaIDFrom(s)
	fmt.Println("root", internal.RootName(s.Location), err)
	walkFn := func(k internal.Kind, prev, in *js.Schema, propName string) error {
		if k == internal.Ref {
			fmt.Printf("%s %s\n", k, in)
		}
		return nil
	}
	return internal.Walk(s, walkFn)
}
