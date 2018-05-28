/**
 * Package for formatting a file as HTML.
 */

package file

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"net/http"
)

var fileTpl *template.Template

const FILE_TEMPLATE = "file.html"

func init() {

	// In DBGp, the line number for the first line is one rather than zero.
	funcMap := template.FuncMap{
		"plusone": func(arg int) int {
			return arg + 1
		},
	}

	fileTpl = template.Must(template.New(FILE_TEMPLATE).Funcs(funcMap).Parse(fileTemplate))
}

/**
 * HTML-ify the given file.
 *
 * Prepare HTML markup for the given file.  The file content is HTML escaped.
 */
func GrabIt(codebase http.Dir, path string) (formattedFile string, err error) {

	fileDesc, err := codebase.Open(path)
	if err != nil {
		return formattedFile, err
	}
	defer fileDesc.Close()

	lines, err := readFileLines(fileDesc)
	formattedFile = formatFile(lines)

	return formattedFile, err
}

/**
 * Prepare a list of lines for the given file.
 */
func grabFileLines(codebase http.Dir, path string) (lines []string, err error) {

	fileDesc, err := codebase.Open(path)
	if err != nil {
		return lines, err
	}
	defer fileDesc.Close()

	lines, err = readFileLines(fileDesc)

	return lines, err
}

/**
 * Read and split a file into its lines.
 */
func readFileLines(in io.Reader) (lines []string, err error) {

	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, err
}

/**
 * Prepare HTML markup for a file's content.
 */
func formatFile(lines []string) string {

	var buffer bytes.Buffer

	fileTpl.ExecuteTemplate(&buffer, FILE_TEMPLATE, lines)

	return buffer.String()
}
