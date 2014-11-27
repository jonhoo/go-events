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
