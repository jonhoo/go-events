package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEqualityDispatch(t *testing.T) {
	Verbosity = 3

	c := Listen(".")
	e := Event{".", nil}

	go Announce(e)

	time.Sleep(5 * time.Millisecond)
	require.Equal(t, e, <-c, "listener should receive announced event")
	require.Empty(t, c, "listener should be notified only once of matching event")
}

func TestNonEqualityDispatch(t *testing.T) {
	Verbosity = 3

	c := Listen("nomatch")
	e := Event{".", nil}

	go Announce(e)

	time.Sleep(5 * time.Millisecond)
	require.Empty(t, c, "listener should not be notified of non-matching event")
}

func TestClose(t *testing.T) {
	Verbosity = 3

	c := Listen(".")
	e := Event{".", nil}
	close(c)
	require.NotPanics(t, func() {
		Announce(e)
	}, "announce should not panic if listener leaves")
}

func TestEventsIndependence(t *testing.T) {
	Verbosity = 3

	ev1 := New()
	ev2 := New()

	c1 := ev1.Listen(".")
	c2 := ev2.Listen(".")
	cGlobal := Listen(".")
	e := Event{".", nil}

	go ev1.Signal(".")
	require.Equal(t, e, <-c1)
	require.Empty(t, c2, "other Events should not be notified")
	require.Empty(t, cGlobal,
		"the global Events should not be notified for created Events")

	go Signal(".")
	require.Equal(t, e, <-cGlobal)
	require.Empty(t, c1, "global signals should not trigger others")
	require.Empty(t, c2, "global signals should not trigger others")
}
