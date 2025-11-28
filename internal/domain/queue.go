package domain

type node struct {
	Coords Object
	Next   *node
}

type queue struct {
	Begin *node
	End   *node
}

func createQueue() *queue {
	return &queue{
		Begin: nil,
		End:   nil,
	}
}

func createNode(coords *Object) *node {
	return &node{
		Coords: *coords,
		Next:   nil,
	}
}

func (q *queue) isEmpty() bool {
	return q.Begin == nil
}

func (q *queue) enqueue(coords *Object) {
	node := createNode(coords)

	if q.End == nil {
		q.Begin = node
		q.End = node
	} else {
		q.End.Next = node
		q.End = node
	}
}

func (q *queue) dequeue() Object {
	if q.isEmpty() {
		return Object{}
	}

	temp := q.Begin
	q.Begin = q.Begin.Next
	if q.Begin == nil {
		q.End = nil
	}

	coords := temp.Coords
	return coords
}
