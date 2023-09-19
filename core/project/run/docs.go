package run

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/tmpl"
)

const (
	docsGenDir    = "docs"
	archGenDir    = "architecture"
	projectConfig = "config/jarvis.json"
	readme        = "README.md"
	docsHeader    = "# Docs"
	serviceHeader = "## Services"
	alignLeft     = ":-"
	alignCentre   = ":-:"
)

//go:embed all:_templates/docs
var docFiles embed.FS

var Fswalkdir = fs.WalkDir

// Generates a table of contents for ./docs/<topic>/README.md files
// and a list of all services in a jarvis project.
// TargetDir should be a jarvis project root directory.
// If check is true, GenerateDocs will instead check the current docs
// and exit with error code 1 if docs are not up to date.
//
// TODO generate graph
func GenerateDocs(p *api_v1.PluginPlaceholders, targetDir string, check bool, t tmpl.TmplWriter) error {
	docsTemplateDir := filepath.Join(goTemplateDir, docsGenDir)
	docsPath := filepath.Join(targetDir, docsGenDir)

	if check {
		tempPath, err := os.MkdirTemp("", "tempdocs*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempPath)

		defer checkDocs(docsPath, tempPath)

		docsPath = tempPath
	}

	if err := tmpl.GenerateFS(docFiles, docsTemplateDir, docsPath, p, false, t); err != nil {
		return fmt.Errorf("regenerating docs file structure: %w", err)
	}

	if err := writeDocsTable(docsPath); err != nil {
		return fmt.Errorf("writing table of contents in %s: %w", docsPath, err)
	}

	servicesPath := filepath.Join(targetDir, p.ServicesDirectory)
	if file.Exists(servicesPath) {
		if err := writeServiceTable(docsPath, servicesPath, p.ServicesDirectory); err != nil {
			return fmt.Errorf("writing table of services in %s: %w", filepath.Join(docsPath, archGenDir), err)
		}
	}

	return nil
}

func writeDocsTable(docsPath string) error {
	docsReadme := filepath.Join(docsPath, readme)
	lines, err := file.ToLines(docsReadme)
	if err != nil {
		return fmt.Errorf("writing readme contents to slice of strings: %w", err)
	}
	table := []string{
		"\n",
		createRow("Topic", "Info"),
		createRow(alignLeft, alignCentre),
	}
	rows, err := fsTableRows(docsPath, 0, "")
	if err != nil {
		return fmt.Errorf("creating table rows for docs directories: %w", err)
	}
	table = append(table, rows...)

	lines = file.InsertIntoLines(lines, docsHeader, table...)
	err = file.FromLines(docsReadme, lines)
	if err != nil {
		return fmt.Errorf("writing table to docs readme: %w", err)
	}

	return nil
}

func writeServiceTable(docsPath, servicesPath, servicesDirectory string) error {
	archReadme := filepath.Join(docsPath, archGenDir, readme)
	lines, err := file.ToLines(archReadme)
	if err != nil {
		return fmt.Errorf("writing readme contents to slice of strings: %w", err)
	}
	table := []string{
		"\n",
		createRow("Namespace", "Service", "Info"),
		createRow(alignLeft, alignLeft, alignCentre),
	}

	linkPrefix := filepath.Join("..", "..", servicesDirectory)
	rows, err := fsTableRows(servicesPath, 1, linkPrefix)
	if err != nil {
		return fmt.Errorf("creating table rows for all services: %w", err)
	}
	table = append(table, rows...)

	lines = file.InsertIntoLines(lines, serviceHeader, table...)
	err = file.FromLines(archReadme, lines)
	if err != nil {
		return fmt.Errorf("writing table to docs readme: %w", err)
	}

	return nil
}

// Returns a table row for each directory in dirPath that contains a README.
// SubDirs specifies the number of directories between the dirPath and the target directories.
// Each subDir will be included in the row for the first README under that directory.
func fsTableRows(dirPath string, subDirs int, linkPrefix string) ([]string, error) {
	rows := []string{}
	directories := make([]string, subDirs)

	err := Fswalkdir(os.DirFS(dirPath), ".",
		func(path string, d fs.DirEntry, err error) error {
			if errors.Is(err, fs.ErrPermission) {
				return fs.SkipDir
			}
			nesting := strings.Count(path, string(os.PathSeparator))
			if !d.IsDir() || nesting > subDirs || path == "." {
				return nil
			}

			readmePath := filepath.Join(dirPath, path, readme)
			linkPath := filepath.Join(linkPrefix, path)

			if nesting < subDirs {
				directories[nesting] = d.Name()
				if file.Exists(readmePath) {
					directories[nesting] = createLink(d.Name(), linkPath)
				}
				return nil
			}

			info, err := infoFromReadme(readmePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fs.SkipDir
				}
				return fmt.Errorf("extracting info from readme in %s: %w", readmePath, err)
			}
			link := createLink(info.header, linkPath)
			row := createRow(append(directories, link, info.desc)...)
			rows = append(rows, row)
			directories = make([]string, subDirs)
			return nil
		})
	return rows, err
}

func createLink(text string, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

func createRow(elems ...string) string {
	joined := strings.Join(elems, " | ")
	return fmt.Sprintf("| %s |", joined)
}

type readmeInfo struct {
	header string
	desc   string
}

// Returns the first header (# characters removed) and subsequent text line from a README.md as a description.
// If no headers are found path base is used as the header.
func infoFromReadme(path string) (*readmeInfo, error) {
	info := &readmeInfo{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			if info.header == "" {
				header := strings.ReplaceAll(scanner.Text(), "#", "")
				info.header = strings.TrimSpace(header)
			}
			continue
		}
		if line != "" {
			info.desc = line
			break
		}
	}

	if info.header == "" {
		info.header = filepath.Base(path)
	}

	return info, scanner.Err()
}

func checkDocs(existingPath string, genPath string) {
	same, err := areDocsSame(existingPath, genPath)
	if err != nil {
		log.Fatalf("error checking docs: %v", err)
	}
	if !same {
		log.Fatal("docs not up to date")
	}
	fmt.Println("docs are up to date")
}

func areDocsSame(path1 string, path2 string) (bool, error) {
	files, err := getDocsTemplateFiles()
	if err != nil {
		return false, fmt.Errorf("retrieving paths for template files: %w", err)
	}

	for _, f := range files {
		f1 := filepath.Join(path1, f)
		f2 := filepath.Join(path2, f)

		same, err := areFilesSame(f1, f2)
		if err != nil {
			return false, fmt.Errorf("comparing equality of %s and %s: %w", f1, f2, err)
		}
		if !same {
			return false, nil
		}
	}

	return true, nil
}

func getDocsTemplateFiles() ([]string, error) {
	files := []string{}
	root := filepath.Join(goTemplateDir, docsGenDir)

	err := fs.WalkDir(docFiles, root, func(fullpath string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrPermission) {
			return fs.SkipDir
		}
		path := filepath.Clean(strings.Replace(fullpath, root, "", 1))[1:]

		if d.Type().IsRegular() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func areFilesSame(file1 string, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	chunkSize := 4000
	b1 := make([]byte, chunkSize)
	b2 := make([]byte, chunkSize)

	for {
		_, err1 := f1.Read(b1)
		if err1 != nil && !errors.Is(err1, io.EOF) {
			return false, err1
		}
		_, err2 := f2.Read(b2)
		if err2 != nil && !errors.Is(err2, io.EOF) {
			return false, err2
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		}
		if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}
