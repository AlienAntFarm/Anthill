package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/alienantfarm/anthive/common"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const assets_go = `
package assets

const (
  {{- range $path, $content:= .Assets }}
  {{ normalize $path }} = ` + "`{{printf \"%s\" $content}}`" + `
  {{- end }}
)

var assets = map[string]string{
  {{- range $path, $_ := .Assets }}
  "{{ $path }}": {{ normalize $path }},
  {{- end }}
}

func Get(assetPath string) string {
  return assets[assetPath]
}
`

func parseArgs() []string {
	flag.Parse()
	paths := flag.Args()

	return paths
}

func moveToAssets() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// if no GOPATH, abort
		common.Error.Fatalf("To generate embeded string, we need a GOPATH")
	}
	assets_path := []string{gopath, "src"}
	assets_path = append(assets_path, strings.Split(common.PROJECT, "/")...)
	assets_path = append(assets_path, "assets")

	err := os.Chdir(path.Join(assets_path...))
	if err != nil {
		common.Error.Fatalf("something bad happened when moving to assets path: %s", err)
	}
}

func normalize(s string) (string, error) {
	var err error = nil
	var translate = func(r rune) rune {
		switch {
		case r == '.' || r == '/' || r == '\\':
			r = '_'
		case r >= 'A' && r <= 'Z':
			r = r + 26
		case r > 127:
			panic(fmt.Sprintf("Non ascii character %s, in %s", r, s))
		}
		return r
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return strings.Map(translate, s), err
}

func main() {
	funcs := template.FuncMap{"normalize": normalize}
	context := struct{ Assets map[string][]byte }{make(map[string][]byte)}

	files := parseArgs()
	moveToAssets()

	for _, f := range files {
		paths, err := filepath.Glob(f)
		if err != nil {
			common.Error.Fatalf(err.Error())
		}
		for _, p := range paths {
			err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || path == "assets.go" {
					return nil
				}
				common.Info.Printf("parsing: %s", path)
				context.Assets[path], err = ioutil.ReadFile(path)
				return err
			})
			if err != nil {
				common.Info.Fatalf("%s when traversing %s", err, p)
			}
		}
	}
	tmpl, err := template.New("assets_go").Funcs(funcs).Parse(assets_go)
	if err != nil {
		common.Error.Fatalf("error when parsing asset template, %s, aborting...", err)
	}
	out, err := os.Create("assets.go")
	if err != nil {
		common.Error.Fatalf("could not generate assets.go file, %s", err)
	}

	err = tmpl.Execute(out, context)
	out.Close()
	if err != nil {
		os.Remove("assets.go")
		common.Error.Fatalf("error when generating the assets.go file, %s, aborting...", err)
	}
}
