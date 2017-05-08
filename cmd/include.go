package main

import (
	"errors"
	"flag"
	"fmt"
	_ "github.com/alienantfarm/anthive/utils" // be sure that viper init is run
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func moveToAssets() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// if no GOPATH, abort
		return errors.New("To generate embeded string, we need a GOPATH")
	}
	assets_path := strings.Split(viper.Get("PROJECT").(string), "/")
	assets_path = append([]string{gopath, "src"}, assets_path...)
	assets_path = append(assets_path, "assets")

	return os.Chdir(path.Join(assets_path...))
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

var rootCmd = &cobra.Command{
	Use:     "include",
	Short:   "parse assets glob under assets dir",
	Example: "include sql/*.sql: will parse all assets in assets/sql dir \n\twhich have .sql as file extension",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		os.Args = os.Args[:1]
		flag.Set("logtostderr", "true")
		flag.Parse()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := generate(args)
		if err != nil {
			glog.Fatalf("%s", err)
		}
	},
}

func generate(files []string) (err error) {
	funcs := template.FuncMap{"normalize": normalize}
	context := struct{ Assets map[string][]byte }{make(map[string][]byte)}

	err = moveToAssets()
	if err != nil {
		return err
	}

	for _, f := range files {
		paths, err := filepath.Glob(f)
		if err != nil {
			glog.Fatalf(err.Error())
		}
		for _, p := range paths {
			err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || path == "assets.go" {
					return nil
				}
				glog.Infof("parsing: %s", path)
				context.Assets[path], err = ioutil.ReadFile(path)
				return err
			})
			if err != nil {
				glog.Fatalf("%s when traversing %s", err, p)
			}
		}
	}
	tmpl, err := template.New("assets_go").Funcs(funcs).Parse(assets_go)
	if err != nil {
		return
	}
	out, err := os.Create("assets.go")
	if err != nil {
		return
	}
	defer out.Close()

	err = tmpl.Execute(out, context)
	if err != nil {
		defer os.Remove("assets.go")
	}
	return
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
