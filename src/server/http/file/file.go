/**
 * Package for formatting a file as HTML.
 */

package file

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

var fileTpl *template.Template

const FILE_TEMPLATE = "file.html"

func init() {

	fileTpl = template.Must(template.New(FILE_TEMPLATE).Parse(fileTemplate))
}

/**
 * HTML-ify the given file.
 *
 * Prepare HTML markup for the given file.  The file content is HTML escaped.
 */
func GrabIt(path string, port int) (formattedFile string, err error) {

	lines, err := grabFileLines(path, port)

	formattedFile = formatFile(lines)

	return formattedFile, err
}

/**
 * Prepare a list of lines for the given file.
 *
 * For now, fetch the file over HTTP.  This should be replaced with a local file
 * operation once we have better understanding of the security implications.
 */
func grabFileLines(path string, port int) (lines []string, err error) {

	fileUrl := fmt.Sprintf("http://0.0.0.0:%d/files/%s", port, path)
	response, err := http.Get(fileUrl)

	if err != nil {
		return lines, err
	}

	defer response.Body.Close()

	lines, err = readFileLines(response.Body)

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
