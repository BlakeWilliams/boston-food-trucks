package bostontrucks

import (
	"time"
)

const DefaultTruckCacheDuration = 12 * time.Hour

// FoodTruckFinder is an interface for finding food trucks.
type TruckFinder interface {
	Trucks() ([]Truck, error)
}

// Manager handles the fetching and caching of food trucks. Caching is implemented to avoid hitting the API too often.
type Manager struct {
	finder    TruckFinder
	trucks    []Truck
	lastFetch time.Time
	// The duration after which the trucks are fetched again.
	CacheDuration time.Duration
}

// Returns a new manager using the given finder to fetch trucks.
func NewManager(finder TruckFinder) *Manager {
	return &Manager{finder: finder, CacheDuration: DefaultTruckCacheDuration}
}

func (tm *Manager) Trucks() ([]Truck, error) {
	if !tm.lastFetch.IsZero() && time.Now().Sub(tm.lastFetch) < tm.CacheDuration {
		return tm.trucks, nil
	}

	trucks, err := tm.finder.Trucks()
	if err != nil {
		return nil, err
	}

	// cache trucks
	tm.trucks = trucks
	tm.lastFetch = time.Now()
	return tm.trucks, nil
}
