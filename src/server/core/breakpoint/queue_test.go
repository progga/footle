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

/**
 * Tests for the delete operation from the breakpoint Queue.
 */
func TestQueueDelete(t *testing.T) {

	var q Queue

	// Push two items and then remove the first one.
	q.push(breakpoint{Filename: "foo.php"})
	q.push(breakpoint{Filename: "bar.php"})

	q.delete(0)

	remaining := q.pop()
	expectedFilename := "bar.php"
	if remaining.Filename != expectedFilename {
		t.Errorf("After deleting the first item, the popped item should be for %s", expectedFilename)
	}
}
