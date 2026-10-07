// Command import-boundaries checks directed package dependencies in the Go
// module. Existing service-to-transport edges remain legacy; new layers are
// checked against the target dependency direction.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const module = "github.com/anlekg/quarel"

type goPackage struct {
	ImportPath   string
	Imports      []string
	TestImports  []string
	XTestImports []string
}

func main() {
	cmd := exec.Command("go", "list", "-test", "-json", "./internal/...")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		fatalf("go list failed: %v\n%s", err, stderr.String())
	}
	dec := json.NewDecoder(&stdout)
	var packages []goPackage
	for {
		var p goPackage
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			fatalf("decode go list output: %v", err)
		}
		packages = append(packages, p)
	}
	if err := validate(packages); err != nil {
		fatalf("%v", err)
	}
	fmt.Printf("import boundaries OK (%d internal packages checked)\n", len(packages))
}

func validate(packages []goPackage) error {
	if !hasPackageRoot(packages, module+"/internal/platform") {
		return fmt.Errorf("no package under internal/platform; boundary check would be vacuous")
	}
	if !hasPackageRoot(packages, module+"/internal/adapters") {
		return fmt.Errorf("no package under internal/adapters; boundary check would be vacuous")
	}
	bad := violations(packages)
	if len(bad) > 0 {
		return fmt.Errorf("forbidden imports:\n%s", strings.Join(bad, "\n"))
	}
	return nil
}

func hasPackageRoot(packages []goPackage, root string) bool {
	for _, p := range packages {
		if p.ImportPath == root || strings.HasPrefix(p.ImportPath, root+"/") {
			return true
		}
	}
	return false
}

func violations(packages []goPackage) []string {
	var bad []string
	for _, p := range packages {
		layer := layerOf(p.ImportPath)
		imports := append(append(append([]string{}, p.Imports...), p.TestImports...), p.XTestImports...)
		for _, imported := range imports {
			if forbidden(layer, imported) {
				bad = append(bad, p.ImportPath+" -> "+imported)
			}
		}
	}
	sort.Strings(bad)
	return bad
}

func layerOf(path string) string {
	for _, layer := range []string{"platform", "application", "adapters", "transport"} {
		root := module + "/internal/" + layer
		if path == root || strings.HasPrefix(path, root+"/") {
			return layer
		}
	}
	return ""
}

func forbidden(from, imported string) bool {
	internal := module + "/internal/"
	if !strings.HasPrefix(imported, internal) {
		return false
	}
	switch from {
	case "platform":
		return !matchesRoot(imported, internal+"platform")
	case "application":
		return !matchesRoot(imported, internal+"application") && !matchesRoot(imported, internal+"platform")
	}
	if from == "adapters" {
		for _, core := range []string{"community", "identity", "httpapi", "realtime", "voice", "backup", "sqlitedb"} {
			if imported == internal+core || strings.HasPrefix(imported, internal+core+"/") {
				return true
			}
		}
		return matchesRoot(imported, internal+"transport")
	}
	if from == "transport" {
		for _, core := range []string{"adapters", "community", "identity", "httpapi", "realtime", "voice", "backup", "sqlitedb"} {
			if matchesRoot(imported, internal+core) {
				return true
			}
		}
	}
	return false
}

func matchesRoot(path, root string) bool { return path == root || strings.HasPrefix(path, root+"/") }

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
