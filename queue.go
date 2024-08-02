package main

type Node struct {
	value Track
	next  *Node //use of linked list to implements queue
}

type TrackQueue struct {
	head *Node
	tail *Node
	size int
}

func (qe *TrackQueue) Enqueue(value Track) {
	newNode := &Node{value: value}
	if qe.size == 0 {
		qe.head = newNode
		qe.tail = newNode
	} else {
		qe.tail.next = newNode
		qe.tail = newNode
	}
	qe.size++
}

func (qe *TrackQueue) Dequeue() (Track, bool) {
	if qe.size == 0 {
		loadItems(qe)
		//return Track{}, false
	}
	value := qe.head.value
	qe.head = qe.head.next
	qe.size--
	return value, true
}

func (qe *TrackQueue) Size() int {
	return qe.size
}

func newMusicQueue() *TrackQueue {
	queue := TrackQueue{}
	loadItems(&queue)
	return &queue
}

func loadItems(queue *TrackQueue) {
	tracks := getTracks()

	for _, track := range tracks {
		queue.Enqueue(track)
	}
}
