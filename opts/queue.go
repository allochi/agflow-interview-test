package opts

// Queue naive queue data structure
type Queue struct {
	values []string
}

// Len returns queue length
func (que *Queue) Len() int {
	return len(que.values)
}

// Push push items into the queue
func (que *Queue) Push(items ...string) {
	que.values = append(que.values, items...)
}

// Pop pop items from the queue
func (que *Queue) Pop() string {
	node := que.values[0]
	que.values = que.values[1:]
	return node
}
