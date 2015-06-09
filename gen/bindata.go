package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const bindataPrefix = "bindata/"

func main() {
	out, err := os.Create("gen_bindata.go")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(out, "package main\n")
	fmt.Fprintf(out, "var binData = map[string]string{\n")

	wf := func(path string, fi os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		path = strings.TrimPrefix(path, bindataPrefix)

		fmt.Fprintf(out, "\"%s\": ", path)

		file, err := os.Open(filepath.Join(bindataPrefix, path))
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Fprintf(out, "\"")

		for _, b := range bytes {
			fmt.Fprintf(out, "\\x%02x", b)
		}

		fmt.Fprintf(out, "\",\n")

		return nil
	}

	err = filepath.Walk(bindataPrefix, wf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(out, "}\n")
}
