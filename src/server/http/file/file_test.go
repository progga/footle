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
		`<div class="lines">
  <pre class="line line__0">foo</pre>
  <pre class="line line__1">bar</pre>
  <pre class="line line__2">buz</pre>
  <pre class="line line__3">qux</pre>
</div>`

	if expectedOutput != formattedFile {
		t.Error("Output mismatch for formatFile().")
	}
}
