package bostontrucks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	manager := NewManager(&MockFinder{})

	trucks, err := manager.Trucks()

	require.NoError(t, err)
	require.Len(t, trucks, 1)
}

type countFinder struct {
	calls int
}

func (cf *countFinder) Trucks() ([]Truck, error) {
	cf.calls++
	return (&MockFinder{}).Trucks()
}

func TestManager_Cache(t *testing.T) {
	finder := &countFinder{}
	manager := NewManager(finder)

	trucks, err := manager.Trucks()
	require.NoError(t, err)
	require.Equal(t, 1, finder.calls)

	trucks, err = manager.Trucks()
	require.NoError(t, err)
	require.Equal(t, 1, finder.calls)

	manager.lastFetch = time.Now().Add(-time.Hour * 48)
	trucks, err = manager.Trucks()
	require.NoError(t, err)
	require.Equal(t, 2, finder.calls)

	require.Len(t, trucks, 1)
}
