package typedmsg

import "fmt"

const queueSize = 100

//DurableQueue is a structure which stores a queue with persistent messages and an ID of a client if he set the durable flag to true;
type DurableQueue struct {
	ClientId      string
	Queue         []string
	DurableTopics []string
}

//NewDurableQueue function which create a *DurableQueue data structure
func NewDurableQueue(clientId string, durableTopics []string) *DurableQueue {
	return &DurableQueue{
		ClientId:      clientId,
		Queue:         []string{},
		DurableTopics: durableTopics,
	}
}

// Dequeue -> method which removes first element from a Queue
func (dq *DurableQueue) Dequeue() {
	element := dq.Queue[0]
	fmt.Println("Dequeued:", element)
	dq.Queue = dq.Queue[1:]
}

// Enqueue -> method which add element to Queue
func (dq *DurableQueue) Enqueue(element string) {
	dq.Queue = append(dq.Queue, element)
}

// Peek -> method to take the first element of the Queue
func (dq *DurableQueue) Peek() string {
	return dq.Queue[0]
}

// IsFull -> method to check if our Queue is full;
func (dq *DurableQueue) IsFull() bool {
	if len(dq.Queue) >= queueSize {
		return true
	}
	return false
}

//setClientId -> method to set clientId for our DurableQueue
func (dq *DurableQueue) setClientId(clientUniqueId string) {
	dq.ClientId = clientUniqueId
}

// IsEmpty -> method to check if our Queue is empty or not
func (dq *DurableQueue) IsEmpty() bool {
	if len(dq.Queue) == 0 {
		return true
	}
	return false
}
