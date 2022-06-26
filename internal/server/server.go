package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/blakewilliams/boston-foodie/pkg/bostontrucks"
	"github.com/blakewilliams/boston-foodie/pkg/static"
	"github.com/blakewilliams/medium/pkg/router"
	"github.com/blakewilliams/medium/pkg/router/rescue"
	"github.com/blakewilliams/medium/pkg/view"
)

type Action struct {
	router.Action
	renderer *view.Renderer
}

func (a *Action) Render(name string, params map[string]any) {
	err := a.renderer.Render(a.ResponseWriter(), name, params)
	if err != nil {
		panic(err)
	}
}

func New() *router.Router[*Action] {
	r := router.New(func(a router.Action) *Action {
		return &Action{Action: a, renderer: createRenderer()}
	})

	r.Use(rescue.Middleware(func(a router.Action, err error) {
		fmt.Println(err)
	}))

	r.Use(static.Middleware(static.Config{
		FileRoot:   "./static",
		PathPrefix: "/static",
	}))

	r.Get("/", func(action *Action) {
		trucks, err := bostontrucks.TrucksByLocation()
		if err != nil {
			panic(err)
		}

		sort.Slice(trucks, func(i, j int) bool {
			return trucks[i].Name < trucks[j].Name
		})

		trucksByLocation := map[string][]bostontrucks.Truck{}

		for _, truck := range trucks {
			location := truck.Neighborhood + " - " + truck.Location

			if _, ok := trucksByLocation[location]; !ok {
				trucksByLocation[location] = make([]bostontrucks.Truck, 0)
			}

			// Backfill seaport location because the food truck page refuses to be consistent
			if truck.Neighborhood == "Seaport" {
				truck.LatLng = bostontrucks.LatLng{Lat: 42.352676, Lng: -71.0452338}
			} else if truck.Neighborhood == "Charlestown" {
				truck.LatLng = bostontrucks.LatLng{Lat: 42.3776592, Lng: -71.0517346}
			}
			trucksByLocation[location] = append(trucksByLocation[location], truck)
		}

		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		today := time.Now().Weekday().String()

		trucksJSON, err := json.Marshal(trucksByLocation)

		if err != nil {
			panic(err)
		}

		action.Render(
			"index.html",
			map[string]any{
				"MapboxAPIKey": os.Getenv("MAPBOX_API_KEY"),
				"Trucks":       trucksByLocation,
				"TrucksJSON":   template.JS(trucksJSON),
				"Days":         days,
				"Today":        today,
			},
		)
	})

	return r
}

func createRenderer() *view.Renderer {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viewRoot := "./views"
	if strings.HasSuffix(cwd, "internal/server") {
		viewRoot = filepath.Join(cwd, "../../views")
	}

	renderer := view.New(viewRoot)
	renderer.Helper("ShortDay", func(dayOfWeek string) string {
		return dayOfWeek[:3]
	})
	err = renderer.AutoRegister()

	if os.Getenv("MEDIUM_ENV") == "development" {
		renderer.HotReload = true

	}

	if err != nil {
		panic(err)
	}

	return renderer
}
