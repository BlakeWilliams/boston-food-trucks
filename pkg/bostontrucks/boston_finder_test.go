package bostontrucks

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed example.html
var bostonPage []byte

func TestTrucksBoston(t *testing.T) {
	trucks, err := parseBostonTrucksHTML(bytes.NewReader(bostonPage))
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

//go:embed example_dewey.html
var deweyPage []byte

func TestTrucksDewey(t *testing.T) {
	trucks, err := parseDeweyTrucksHTML(bytes.NewReader(deweyPage))
	require.NoError(t, err)

	require.Len(t, trucks, 26)

	var bonMe *Truck
	for _, truck := range trucks {
		if truck.Name == "Bon Me" && truck.Location == "Dewey Square" {
			bonMe = &truck
			break
		}
	}

	require.NotNil(t, bonMe)
	require.Equal(t, bonMe.Name, "Bon Me")
	require.Len(t, bonMe.Schedule, 4)
	require.Equal(t, bonMe.Schedule["Monday"], "N/A")
	require.Equal(t, bonMe.Schedule["Tuesday"], "N/A")
	require.Equal(t, bonMe.Schedule["Wednesday"], "N/A")
	require.Equal(t, bonMe.Schedule["Thursday"], "N/A")
}
