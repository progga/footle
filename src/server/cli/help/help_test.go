/**
 * @file
 * Unit tests for the Help class.
 */

package help

import "testing"

/**
 * Tests for Help.add().
 */
func TestAdd(t *testing.T) {

	texts := []helptext{
		helptext{[]string{"foo", "bar"}, "Sample helptext for foo."},
		helptext{[]string{"baz", "qux"}, "Sample helptext for baz."},
		helptext{[]string{"tar"}, "Sample helptext for tar."},
	}

	var testObj Help
	testObj.prepare(texts, []helptext{}, []helptext{})

	if len(texts) != len(testObj.cmdList) {
		t.Errorf("Failed to pick all command list.")
	}

	helptextForQux := (*testObj.cmdHelptextMapping["qux"]).explanation
	expected := texts[1].explanation
	if helptextForQux != expected {
		t.Errorf("Got wrong helptext for Qux: %s", helptextForQux)
	}

	moreTexts := []helptext{
		helptext{[]string{"foo2", "bar2"}, "Sample helptext for foo2."},
	}
	testObj.add(moreTexts)

	helptextForTar := (*testObj.cmdHelptextMapping["tar"]).explanation
	expected = texts[2].explanation
	if helptextForTar != expected {
		t.Errorf("Got wrong helptext for Tar: %s", helptextForTar)
	}

	helptextForFoo2 := (*testObj.cmdHelptextMapping["foo2"]).explanation
	expected = moreTexts[0].explanation
	if helptextForFoo2 != expected {
		t.Errorf("Got wrong helptext for Foo2: %s", helptextForFoo2)
	}
}

/**
 * Tests for Help.forAll().
 */
func TestForAll(t *testing.T) {

	texts := []helptext{
		helptext{[]string{"foo", "bar"}, "Sample helptext for foo."},
		helptext{[]string{"baz", "qux"}, "Sample helptext for baz."},
		helptext{[]string{"tar"}, "Sample helptext for tar."},
	}

	var helpObj Help
	helpObj.prepare(texts, []helptext{}, []helptext{})
	result := helpObj.forAll()

	expected := "foo, bar\nbaz, qux\ntar\n"
	if result != expected {
		t.Errorf("Command list does not meet expectation. %s given.", result)
	}
}
