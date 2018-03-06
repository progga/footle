/**
 * A queue data structure for breakpoints.
 */
package breakpoint

import "errors"

type Queue []breakpoint

/**
 * Push items to the bottom of the Queue.
 */
func (q *Queue) push(b breakpoint) {

	*q = append(*q, b)
}

/**
 * Pop an item from the top of the Queue.
 */
func (q *Queue) pop() (b breakpoint) {

	queueLength := len(*q)

	if queueLength == 0 {
		return b
	}

	b = (*q)[0]

	*q = (*q)[1:]

	return b
}

/**
 * Delete the item for the given array index.
 */
func (q *Queue) delete(index int) (err error) {

	if index < 0 || index >= len(*q) {
		err = errors.New("Array index out of bound.")
		return err
	}

	*q = append((*q)[:index], (*q)[index+1:]...)

	return err
}
