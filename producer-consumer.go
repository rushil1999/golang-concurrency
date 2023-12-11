/*
The producer-consumer problem is a synchronization problem that occurs when more than one process shares a resource with limited resources.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Constant to terminate program when more than 2 events fail
const FAILURE_THRESHOLD = 2
const EVENT_CHANNEL_BUFFER_SIZE = 5 // This determines the capacity of a producer to produce and persist evnts

type Event struct {
	eventId int
	message string
	eventDelay int
}

type Producer struct {
	dataChannel chan Event
	quitChannel chan bool
	failedEvents int
	totalEvents int
	successEvents int
}

type Consumer struct {
	eventsConsumed int 
}

type producerInterface interface {
	generateEvent(id int) Event
	generateEvents(eventChannel chan<- Event) 
}

type consumerInterface interface {
	consumeEvents(eventChannel <-chan Event, quitChannel <-chan bool)
}

// This function will consume events from the eventChannel
func (c *Consumer) consumeEvents(eventChannel <-chan Event, quite <-chan bool) {
		//This loop will consume the channel
		for {
			select {
			case event := <- eventChannel :
				fmt.Printf("Consuming events from channel with event id #%d\n", event.eventId)
				time.Sleep(time.Duration(10*event.eventDelay)*time.Second)
				fmt.Printf("Consumed event #%d\n", event.eventId)
				c.eventsConsumed += 1
			default:
				if c.eventsConsumed > 0 {
					fmt.Printf("Total events consumed #%d\n", c.eventsConsumed)
					return 
				} 
			}
		}
}


//This function will take time to return the generated event
// Kind of like an async call to a third party service
func (p *Producer) generateEvent(id int) Event { 
	eventId := id+1
	randomDelay := rand.Intn(5) // Random delay to generate an event
	msg := fmt.Sprintf("Event generated with event id #%d", eventId)
	newEvent := Event{
		eventId: eventId,
		message: msg,
		eventDelay: randomDelay,
	} 
	time.Sleep(time.Duration(randomDelay)* time.Second)
	return newEvent
}

//Once the event is received this function will push the events to channel
// Events will be pushed to the channel in the order they receive
func (p *Producer) generateEvents(eventChannel chan<- Event, quitChannel chan<- bool) {
	currentId := 0
	for {
		newEvent := p.generateEvent(currentId)
		p.totalEvents += 1
		if newEvent.eventDelay >= 4{ // Random condition to make some events fail to eventually terminate the program
			p.failedEvents += 1
			if p.failedEvents > FAILURE_THRESHOLD {
				quitChannel <- true // Passsing a value to a separate channel to ensure things terminte smoothly
				close(quitChannel)
				close(eventChannel)
				fmt.Printf("Terminating event generation.")
				fmt.Printf("Total Events Cosumed %d\n", p.totalEvents)
				fmt.Printf("Total Events Failed %d\n", p.failedEvents)
				return
			}
		} else {
			
			select {
			case eventChannel <- newEvent : //This is to ensure non-blocking send to the events channel
				fmt.Printf("Generated and Pushed event with event id #%d\n", newEvent.eventId)
			}
		}
		currentId = newEvent.eventId
	}
}



func main() {

	producer := &Producer{} //Creating Producer instance

	consumer := &Consumer{} //Creating Consumer instance

	eventChannel := make(chan Event, EVENT_CHANNEL_BUFFER_SIZE) // Creating events channel with capcaity
	quitChannel := make(chan bool) // Creating quite channel 

	// Producer that will run in the background and generate events.
	go producer.generateEvents(eventChannel, quitChannel)
	go consumer.consumeEvents(eventChannel, quitChannel)


	<- quitChannel // This is to make sure the main function does not terminate unless other goroutines are implemented

}