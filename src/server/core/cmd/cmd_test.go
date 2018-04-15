/**
 * @file
 * Tests for validating internal Footle commands.
 */

package cmd

import "testing"

/**
 * Tests for Validate().
 */
func TestValidate(t *testing.T) {

	// Pass cases.
	if _, err := Validate("continue", []string{}); err != nil {
		t.Error("Misidentified valid continue command.")
	}

	if _, err := Validate("update_source", []string{"foo.php"}); err != nil {
		t.Error("Misidentified valid update_source command.")
	}

	// Fail cases.
	if _, err := Validate("continue", []string{"12"}); err == nil {
		t.Error("Failed to spot invalid continue command.")
	}

	if _, err := Validate("update_source", []string{}); err == nil {
		t.Error("Failed to spot invalid update_source command.")
	}

	if _, err := Validate("update_source", []string{"foo.php", "bar.php"}); err == nil {
		t.Error("Failed to spot invalid update_source command.")
	}

	if _, err := Validate("foo", []string{}); err == nil {
		t.Error("Failed to spot invalid command.")
	}
}
