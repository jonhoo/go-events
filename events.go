// Package events provides an event notification and subscription system.
package events

import (
	"fmt"
	"strings"
)

// Event represents a single occurance of some event in the underlying system.
type Event struct {
	// Tag is a string describing the nature of the event. Listeners can
	// subscribe to prefixes of these tags.
	Tag string
	// Data allows arbitrary additional information to be dispatched along
	// with the event to any interested listeners.
	Data interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("<Event %s with data %v>", e.Tag, e.Data)
}

/*
Verbosity indicates how verbose the events library should be.

  0 = no output (default)
  1 = print event occurances
  2 = print event subscriptions
  3 = print event notifications
*/
var Verbosity = 0

type listener struct {
	prefix string
	sendon chan<- Event
}

var announcements = make(chan Event)
var listeners = make(chan listener)

// dispatch will notify the listener c about Event e iff
// e.Tag.HasPrefix(c.prefix). If the listener has left (i.e. closed the
// channel), the function returns false, otherwise it returns true.
func dispatch(c listener, e Event) (present bool) {
	defer func() {
		if r := recover(); r != nil {
			// The given client is no longer listening
			present = false
		}
	}()

	present = true
	if strings.HasPrefix(e.Tag, c.prefix) {
		if Verbosity > 2 {
			fmt.Printf("events: event '%s' dispatched to client listening for '%s'\n", e.Tag, c.prefix)
		}
		c.sendon <- e
	}
	return
}

// init will set up the event dispatching loop. This function *must* be called
// for event dispatching to become operational, but this is handled by the Go
// package loader.
func init() {
	go func() {
		var clients []listener
		for {
			select {
			case c := <-listeners:
				if Verbosity > 1 {
					fmt.Printf("events: client registered interest in prefix '%s'\n", c.prefix)
				}
				clients = append(clients, c)
			case e := <-announcements:
				if Verbosity > 0 {
					fmt.Printf("events: event '%s' occured; data: %v\n", e.Tag, e.Data)
				}

				for _, c := range clients {
					// TODO(jon): By not checking the
					// return value here, we are continuing
					// to check clients that have closed
					// their channels for matches, which
					// could end up being slow. However,
					// checking the return value means
					// either we have to not execute it
					// asynchronously (potentially causing
					// deadlock), or we have to have a list
					// that can be concurrently deleted
					// from; neither of which are things we
					// want to do just now. Patches are
					// welcome.
					go dispatch(c, e)
				}
			}
		}
	}()
}

// Listen will register interest in Events whose Tag have the indicated prefix.
// The returned channel will receive any future matching Event.
//
// To stop listening, simply close the channel. Attempting to send on the
// returned channel yields undefined behaviour. It is only writeable to allow
// closing.
func Listen(prefix string) chan Event {
	io := make(chan Event)
	listeners <- listener{prefix, io}
	return io
}

// Announce will register the occurance of the given Event e, and pass it along
// for dispatch to interested listeners.
func Announce(e Event) {
	announcements <- e
}

// Signal is a shortcut for announcing an event with only a Tag, and no Data.
func Signal(tag string) {
	Announce(Event{tag, nil})
}
