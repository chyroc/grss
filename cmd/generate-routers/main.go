package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func main() {
	pkgList := []*FeedPkg{}
	err := walk("internal/routers", func(path string, file fs.FileInfo) error {
		if file.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if path == "internal/routers/routers.go" {
			return nil
		}
		pkg := readInfo(path)
		pkgList = append(pkgList, pkg)
		return nil
	})

	sort.Slice(pkgList, func(i, j int) bool {
		return pkgList[i].RoutePath < pkgList[j].RoutePath
	})

	text := generateRouters(pkgList)
	err = ioutil.WriteFile("internal/routers/routers.go", []byte(text), 0o666)
	if err != nil {
		panic(err)
	}
}

func walk(dir string, f func(path string, file fs.FileInfo) error) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range fs {
		err = f(dir+"/"+v.Name(), v)
		if err != nil {
			return err
		}
		if v.IsDir() {
			err = walk(dir+"/"+v.Name(), f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type FeedPkg struct {
	RoutePath   string
	PackagePath string
	PackageName string
	Args        []string
}

func readInfo(path string) *FeedPkg {
	// internal/routers/woshipm/latest/feed.go -> internal/routers/woshipm/latest
	// internal/routers/woshipm/latest -> github.com/chyroc/grss/internal/routers/woshipm/latest
	// internal/routers/woshipm/latest -> woshipm/latest
	// woshipm/latest -> woshipm_latest
	args := []string{}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	route := dir[len("internal/routers/"):]
	packagePath := "github.com/chyroc/grss/" + dir
	packageName := strings.ReplaceAll(route, "/", "_")
	if base != "feed.go" {
		// feed_:id.go -> [':id']
		base = base[len("feed_") : len(base)-len(".go")]
		baseArgs := strings.Split(base, "_")
		route = strings.Join(append(strings.Split(route, "/"), baseArgs...), "/")

		bs, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		for _, vv := range strings.Split(string(bs), "\n") {
			vv = strings.TrimSpace(vv)
			if strings.HasPrefix(vv, "// args:") {
				vv = strings.TrimSpace(vv[len("// args:"):])
				args = append(args, vv)
			}
		}
	}
	return &FeedPkg{
		RoutePath:   "/" + route,
		PackagePath: packagePath,
		PackageName: packageName,
		Args:        args,
	}
}

func generateRouters(pkgList []*FeedPkg) string {
	buf := new(bytes.Buffer)
	t, err := template.New("").Parse(routersTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(buf, pkgList)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

var routersTemplate = `package routers

import (
	"github.com/chyroc/grss/internal/fetch"
{{ range . }}"{{ .PackagePath }}"
{{ end }}
)

func init() {
{{ range . }}Register("{{ .RoutePath }}", {{ .PackageName }}.New,
{{ if .Args }}{{ range .Args}} {{ . }},
{{ end }} {{ end }}
)
{{ end }}
}

var routers []Router

type Router struct {
	Path   string
	Source func(map[string]string) (*fetch.Source, error)
	Args   []map[string]string
}

func Register(path string, source func(map[string]string) (*fetch.Source, error), args ...map[string]string) {
	if len(args) == 0 {
		args = []map[string]string{nil}
	}
	routers = append(routers, Router{
		Path:   path,
		Source: source,
		Args:   args,
	})
}

func Get() []Router {
	return routers
}
`
