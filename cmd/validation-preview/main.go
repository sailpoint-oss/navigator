// Command validation-preview walks YAML OpenAPI fixtures and writes a Markdown
// report of navigator validation diagnostics (for PR comment previews).
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sailpoint-oss/navigator"
)

const maxIssuesPerFile = 40

func main() {
	root := flag.String("root", ".", "repository root (contains testdata/)")
	outPath := flag.String("out", "", "write Markdown here (default: stdout)")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "root: %v\n", err)
		os.Exit(1)
	}

	var out io.Writer = os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "out: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if err := writeReport(out, absRoot); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func writeReport(w io.Writer, repoRoot string) error {
	testdata := filepath.Join(repoRoot, "testdata")
	var paths []string
	err := filepath.WalkDir(testdata, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no YAML files under %s", testdata)
	}
	sort.Strings(paths)

	rel := func(abs string) string {
		r, err := filepath.Rel(repoRoot, abs)
		if err != nil {
			return abs
		}
		return filepath.ToSlash(r)
	}

	fmt.Fprintf(w, "## Navigator validation preview\n\n")
	fmt.Fprintf(w, "_OpenAPI fixtures under `testdata/` parsed with `ParseContent` (tree-sitter YAML). ")
	fmt.Fprintf(w, "Up to %d issues per file; severities: error, warning, info._\n\n", maxIssuesPerFile)

	for _, abs := range paths {
		data, err := os.ReadFile(abs)
		if err != nil {
			return err
		}
		uri := "file:///" + filepath.ToSlash(abs)
		idx := navigator.ParseContent(data, uri)
		rp := rel(abs)
		fmt.Fprintf(w, "### `%s`\n\n", rp)
		if idx == nil {
			fmt.Fprintf(w, "_Parse failed or empty input._\n\n")
			continue
		}
		issues := idx.Issues
		if len(issues) == 0 {
			fmt.Fprintf(w, "_No diagnostics._\n\n")
			continue
		}
		writeIssueTable(w, issues)
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "---\n\n")
	fmt.Fprintf(w, "_This comment is created or updated on each push to the PR (workflow `%s`)._\n",
		"validation-preview-comment.yml")
	return nil
}

func writeIssueTable(w io.Writer, issues []navigator.Issue) {
	trunc := false
	if len(issues) > maxIssuesPerFile {
		issues = issues[:maxIssuesPerFile]
		trunc = true
	}
	fmt.Fprintf(w, "| Severity | Category | Code | Pointer | Message |\n")
	fmt.Fprintf(w, "| --- | --- | --- | --- | --- |\n")
	for _, iss := range issues {
		fmt.Fprintf(w, "| %s | %s | `%s` | %s | %s |\n",
			severityLabel(iss.Severity),
			categoryLabel(iss.Category),
			escapePipe(iss.Code),
			escapePipe(pointerCell(iss.Pointer)),
			escapePipe(oneLine(iss.Message)),
		)
	}
	if trunc {
		fmt.Fprintf(w, "\n_Showing first %d issues._\n", maxIssuesPerFile)
	}
}

func pointerCell(p string) string {
	if p == "" {
		return "—"
	}
	return "`" + p + "`"
}

func escapePipe(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "\n", " ")), " ")
}

func severityLabel(s navigator.Severity) string {
	switch s {
	case navigator.SeverityError:
		return "**error**"
	case navigator.SeverityWarning:
		return "warning"
	case navigator.SeverityInfo:
		return "info"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

func categoryLabel(c navigator.IssueCategory) string {
	switch c {
	case navigator.CategorySyntax:
		return "syntax"
	case navigator.CategoryStructural:
		return "structural"
	case navigator.CategorySchema:
		return "schema-shape"
	case navigator.CategoryMeta:
		return "meta"
	default:
		return fmt.Sprintf("cat-%d", int(c))
	}
}
