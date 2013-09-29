package main

import (
	"flag"
	"github.com/axiom/rgbtxt"
	"io"
	"os"
	"text/template"
)

var (
	file = flag.String("file", "", "Filename, will use stdin if empty")
)

const goT = `
// Automatically generated

package rgbtxt

import (
	"image/color"
)

var (
{{range .}}	{{.Go}}{{end}})
`

func main() {
	flag.Parse()

	var input io.Reader
	if *file == "" {
		input = os.Stdin
	} else {
		inputF, err := os.Open(*file)
		if err != nil {
			panic(err)
		}
		input = inputF
	}

	seen := make(map[string]bool)
	colorPairs := make(chan rgbtxt.ColorPair, 100)
	dedupPairs := make(chan rgbtxt.ColorPair, 100)
	go func() {
		if err := rgbtxt.ParseLinesChan(input, colorPairs); err != nil {
			panic(err)
		}
	}()

	go func() {
		for pair := range colorPairs {
			if _, ok := seen[pair.Name]; !ok {
				dedupPairs <- pair
			}
			seen[pair.Name] = true
		}
		close(dedupPairs)
	}()

	tmpl, err := template.New("template").Parse(goT)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, dedupPairs)
	if err != nil {
		panic(err)
	}
}
