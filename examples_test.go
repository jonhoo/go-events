package events

import "fmt"

// This example shows the simples example of how to listen for, and announce,
// simple tagged events.
func ExampleSignal() {
	Verbosity = 0

	chn := Listen("example.hello")

	go func() {
		Signal("example.hello")
	}()

	e := <-chn
	fmt.Println(e.Tag)

	// Output:
	// example.hello
}

// This example demonstrates how tag prefix matching works.
func ExampleListen() {
	Verbosity = 0

	chn := Listen("example.")

	go func() {
		Signal("example.hello.world")
		Signal("example.hello.aliens")
	}()

	i := 0
	for e := range chn {
		fmt.Println(e.Tag)

		// Avoid listening forever
		i++
		if i == 2 {
			break
		}
	}

	// Output:
	// example.hello.world
	// example.hello.aliens
}

// This example shows how to include extra information with an Event beyond the
// Tag, and how this information can be extracted by the receiver.
func ExampleAnnounce() {
	Verbosity = 0

	chn := Listen("example.")

	go func() {
		Announce(Event{"example.hello", "world"})
		Announce(Event{"example.answer", 42})
	}()

	i := 0
	for e := range chn {
		fmt.Println(e.Tag, e.Data)

		// Avoid listening forever
		i++
		if i == 2 {
			break
		}
	}

	// Output:
	// example.hello world
	// example.answer 42
}
