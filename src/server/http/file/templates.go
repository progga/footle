/**
 * Declares HTML templates for formatting a file as HTML.
 *
 * Ideally, templates should be in their own file for ease of editing.  But that
 * makes it difficult to locate them from Unit tests.  This is because the path
 * of a template is determined *relative* to the executable and the path of the
 * executable differs during tests. So we are keeping the templates in code.
 */

package file

var fileTemplate string = `<table class="lines">
  {{- range $key, $value := . }}
  <tr class="line line__{{ $key }}">
    <td class="line__number">{{ $key }}</td>
    <td class="line__code">{{ . }}</td>
  </tr>
  {{- end }}
</table>`
