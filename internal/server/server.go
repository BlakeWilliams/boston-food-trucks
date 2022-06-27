package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/blakewilliams/boston-foodie/pkg/bostontrucks"
	"github.com/blakewilliams/boston-foodie/pkg/static"
	"github.com/blakewilliams/medium/pkg/router"
	"github.com/blakewilliams/medium/pkg/router/rescue"
	"github.com/blakewilliams/medium/pkg/tell"
	"github.com/blakewilliams/medium/pkg/view"
)

const (
	EnvProd        = "prod"
	EnvDevelopment = "development"
	EnvTest        = "test"
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

func New(env string, logger *log.Logger) *router.Router[*Action] {
	instrumenter := tell.New()
	instrumenter.Subscribe("router.ServeHTTP", func(e tell.Event) {
		req := e.Payload["request"].(http.Request)
		logger.Printf("%s %s - %v\n", req.Method, req.URL.Path, time.Now().Sub(e.Start))
	})

	r := router.New(func(a router.Action) *Action {
		return &Action{Action: a, renderer: createRenderer(env)}
	})
	r.Notifier = instrumenter

	r.Use(rescue.Middleware(func(a router.Action, err error) {
		fmt.Println(err)
		a.Write([]byte("Oops, Something went wrong!"))
	}))

	r.Use(static.Middleware(static.Config{
		FileRoot:   "./static",
		PathPrefix: "/static",
	}))

	var finder bostontrucks.TruckFinder = &bostontrucks.BostonFinder{}
	if env == EnvTest {
		finder = &bostontrucks.MockFinder{}
	}

	truckManager := bostontrucks.NewManager(finder)

	r.Get("/status", func(action *Action) {
		action.Write([]byte("OK"))
	})

	r.Get("/", func(action *Action) {
		trucks, err := truckManager.Trucks()
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
				truck.Lat = 42.352676
				truck.Lng = -71.0452338
			} else if truck.Neighborhood == "Charlestown" {
				truck.Lat = 42.3776592
				truck.Lng = -71.0517346
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

func createRenderer(env string) *view.Renderer {
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

	if env == EnvDevelopment {
		renderer.HotReload = true
	}

	if err != nil {
		panic(err)
	}

	return renderer
}
