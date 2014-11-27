package events_test

import (
	"testing"
	"time"

	events "."
	"github.com/stretchr/testify/require"
)

func TestEqualityDispatch(t *testing.T) {
	c := events.Listen(".")
	e := events.Event{".", nil}

	go events.Announce(e)

	time.Sleep(5 * time.Millisecond)
	require.Equal(t, e, <-c, "listener should receive announced event")
	require.Empty(t, c, "listener should be notified only once of matching event")
}

func TestNonEqualityDispatch(t *testing.T) {
	c := events.Listen("nomatch")
	e := events.Event{".", nil}

	go events.Announce(e)

	time.Sleep(5 * time.Millisecond)
	require.Empty(t, c, "listener should not be notified of non-matching event")
}

func TestClose(t *testing.T) {
	c := events.Listen(".")
	e := events.Event{".", nil}
	close(c)
	require.NotPanics(t, func() {
		events.Announce(e)
	}, "announce should not panic if listener leaves")
}
