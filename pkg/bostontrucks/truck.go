package bostontrucks

type LatLng = struct {
	Lat float64
	Lng float64
}

type Truck struct {
	Name         string
	Neighborhood string
	Location     string
	Schedule     map[string]string
	LatLng       LatLng
	URL          string
}
