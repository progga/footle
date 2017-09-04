/**
 * Tests for formatting a file as HTML.
 */

package file

import "testing"

/**
 * Tests for formatFile().
 */
func TestFormatFile(t *testing.T) {

	lines := []string{
		"foo",
		"bar",
		"buz",
		"qux",
	}

	formattedFile := formatFile(lines)

	expectedOutput :=
		`<table class="lines">
  <tr class="line line__1">
    <td class="line__number">1</td>
    <td class="line__code">foo</td>
  </tr>
  <tr class="line line__2">
    <td class="line__number">2</td>
    <td class="line__code">bar</td>
  </tr>
  <tr class="line line__3">
    <td class="line__number">3</td>
    <td class="line__code">buz</td>
  </tr>
  <tr class="line line__4">
    <td class="line__number">4</td>
    <td class="line__code">qux</td>
  </tr>
</table>`

	if expectedOutput != formattedFile {
		t.Error("Output mismatch for formatFile().")
	}
}
