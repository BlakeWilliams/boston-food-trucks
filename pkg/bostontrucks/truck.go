package bostontrucks

type Truck struct {
	Name         string
	Neighborhood string
	Location     string
	Schedule     map[string]string
	Lat          float64
	Lng          float64
	URL          string
}
