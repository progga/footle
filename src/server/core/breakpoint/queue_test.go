package breakpoint

import "testing"

/**
 * Tests for the breakpoint Queue data structure.
 */
func TestQueue(t *testing.T) {

	var q Queue

	// Tests for the push operation.
	q.push(breakpoint{Filename: "foo.php"})
	q.push(breakpoint{Filename: "bar.php"})

	// Now verify the items in the queue.
	expectedFilename := "foo.php"
	if q[0].Filename != expectedFilename {
		t.Errorf("First queued item should be for %s", expectedFilename)
	}

	expectedFilename = "bar.php"
	if q[1].Filename != expectedFilename {
		t.Errorf("Second queued item should be for %s", expectedFilename)
	}

	// Tests for the pop operation.
	breakpointRecord := q.pop()
	expectedFilename = "foo.php"
	if breakpointRecord.Filename != expectedFilename {
		t.Errorf("Popped item should be for %s", expectedFilename)
	}

	// Verify the remaining items in the queue.
	expectedFilename = "bar.php"
	if q[0].Filename != expectedFilename {
		t.Errorf("After the first pop operation, the first queued item should be for %s", expectedFilename)
	}

	breakpointRecord = q.pop()
	expectedFilename = "bar.php"
	if breakpointRecord.Filename != expectedFilename {
		t.Errorf("Popped item should be for %s", expectedFilename)
	}

	// Queue should be empty since all items have been popped.
	expectedQueueLength := 0
	if len(q) != expectedQueueLength {
		t.Error("Queue should be empty.")
	}
}
