package file

import (
	"github.com/bmatcuk/doublestar/v3"
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
	"strings"
)

// hasMeta reports whether path contains any magic glob characters.
func hasMeta(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

// hasSuperMeta reports whether path contains any super magic glob characters (**).
func hasSuperMeta(path string) bool {
	return strings.Contains(path, "**")
}

type GlobPath struct {
	path         string
	hasMeta      bool
	HasSuperMeta bool
	rootGlob     string
	g            glob.Glob
}

func Compile(path string) (*GlobPath, error) {
	out := GlobPath{
		hasMeta:      hasMeta(path),
		HasSuperMeta: hasSuperMeta(path),
		path:         filepath.FromSlash(path),
	}

	// if there are no glob meta characters in the path, don't bother compiling
	// a glob object
	if !out.hasMeta || !out.HasSuperMeta {
		return &out, nil
	}

	// find the root elements of the object path, the entry point for recursion
	// when you have a super-meta in your path (which are :
	// glob(/your/expression/until/first/star/of/super-meta))
	out.rootGlob = path[:strings.Index(path, "**")+1]
	var err error
	if out.g, err = glob.Compile(path, os.PathSeparator); err != nil {
		return nil, err
	}
	return &out, nil
}

func (g *GlobPath) Match() []string {
	// This string replacement is for backwards compatibility support
	// The original implementation allowed **.txt but the double star package requires **/**.txt
	g.path = strings.ReplaceAll(g.path, "**/**", "**")
	g.path = strings.ReplaceAll(g.path, "**", "**/**")

	files, _ := doublestar.Glob(g.path)
	return files
}
