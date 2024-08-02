package main

var UpdateChan chan Track

type Node struct {
	value Track
	next  *Node //use of linked list to implements queue
}

type TrackQueue struct {
	head        *Node
	tail        *Node
	size        int
	UpdateChan  chan Track
	StopQueuing chan bool
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

func (qe *TrackQueue) SyncUpdate() {
	for {
		select {
		case track := <-qe.UpdateChan:
			{
				qe.Enqueue(track)
			}
		case <-qe.StopQueuing:
			{
				return
			}
		}
	}
}

func newMusicQueue() *TrackQueue {
	queue := TrackQueue{
		UpdateChan: UpdateChan,
	}
	loadItems(&queue)
	go queue.SyncUpdate()
	return &queue
}

func loadItems(queue *TrackQueue) {
	tracks := getTracks()

	for _, track := range tracks {
		queue.Enqueue(track)
	}
}
