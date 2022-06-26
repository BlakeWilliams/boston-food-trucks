package bostontrucks

type MockFinder struct{}

func (mf *MockFinder) Trucks() ([]Truck, error) {
	truck := Truck{
		Name:         "Chicken Rice Guys",
		Neighborhood: "Downtown",
		Location:     "Winter Street",
		Schedule:     map[string]string{"Monday": "11 - 3 p.m."},
		LatLng:       LatLng{Lat: 42.35, Lng: 71.05},
		URL:          "",
	}

	return []Truck{truck}, nil
}

var _ TruckFinder = &MockFinder{}
