package bostontrucks

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed example.html
var examplePage []byte

func TestTrucksByLocation(t *testing.T) {
	trucks, err := parseBostonTrucksHTML(bytes.NewReader(examplePage))
	require.NoError(t, err)

	require.Len(t, trucks, 19)

	var chickenRiceGuys *Truck
	for _, truck := range trucks {
		if truck.Name == "Chicken and Rice Guys" && truck.Neighborhood == "Seaport" {
			chickenRiceGuys = &truck
			break
		}
	}

	require.NotNil(t, chickenRiceGuys)
	fmt.Println(chickenRiceGuys)
	require.Equal(t, chickenRiceGuys.Name, "Chicken and Rice Guys")
	require.Len(t, chickenRiceGuys.Schedule, 1)
	require.Equal(t, chickenRiceGuys.Schedule["Monday"], "11 - 3 p.m.")
}
