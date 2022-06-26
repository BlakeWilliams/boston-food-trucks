package bostontrucks

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ScheduleURL = "https://www.boston.gov/departments/small-business-development/city-boston-food-trucks-schedule"

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

var latLngRegexp = regexp.MustCompile(`maps/@(-?[0-9.]+),(-?[0-9.]+),`)

// There's sometimes an invisible character instead of a space, so we need to
// account for it
var truckDetailsRegexp = regexp.MustCompile(`(.*?),\s+(.*?):[\s*| ]*(.*)`)

func TrucksByLocation() ([]Truck, error) {
	res, err := http.Get("https://www.boston.gov/departments/small-business-development/city-boston-food-trucks-schedule")

	if err != nil {
		return nil, fmt.Errorf("could not fetch trucks: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch trucks, received non-200 response code %d", res.StatusCode)
	}

	return ParseTrucksHTML(res.Body)
}

func ParseTrucksHTML(body io.Reader) ([]Truck, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("could not parse foodtruck page: %w", err)
	}

	trucks := make([]Truck, 0)

	doc.Find("#locations-by-neighborhood + .paragraphs-item-drawers .section-drawers .entity.entity-paragraphs-item").Each(func(i int, s *goquery.Selection) {
		neighborhood := s.Find(".dr-t .field-name-field-title").Text()
		neighborhood = strings.TrimSpace(neighborhood)
		if neighborhood == "" {
			return
		}

		s.Find(".dr-c h4").Each(func(i int, location *goquery.Selection) {
			href, ok := location.Find("a").Attr("href")
			if !ok {
				return
			}

			latLng, err := parseLatLng(href)
			if err != nil {
				return
			}

			truckBuilder := make(map[string]*Truck)

			location.Next().Find("li").Each(func(i int, s *goquery.Selection) {
				truckDetails := strings.TrimSpace(s.Text())
				url, _ := s.Find("a").Attr("href")

				matches := truckDetailsRegexp.FindStringSubmatch(truckDetails)
				if len(matches) != 4 {
					return
				}
				dayOfWeek := strings.TrimSuffix(matches[1], "s")
				hours := matches[2]
				name := matches[3]

				if _, ok := truckBuilder[name]; !ok {
					truckBuilder[name] = &Truck{
						Neighborhood: neighborhood,
						LatLng:       latLng,
						Location:     strings.Replace(location.Text(), "StreeT", "Street", 1), // fix for mistake on food truck site
						Name:         name,
						Schedule:     make(map[string]string),
						URL:          url,
					}
				}

				truck := truckBuilder[name]
				truck.Schedule[dayOfWeek] = hours
			})

			for _, truck := range truckBuilder {
				trucks = append(trucks, *truck)
			}
		})

	})

	if len(trucks) == 0 {
		return trucks, fmt.Errorf("no trucks could be parsed")
	}

	return trucks, nil
}

func parseLatLng(mapURL string) (LatLng, error) {
	match := latLngRegexp.FindStringSubmatch(mapURL)

	if len(match) != 3 {
		return LatLng{}, nil
	}

	lat, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return LatLng{}, fmt.Errorf("could not parse %s: %w", match[1], err)
	}
	lng, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return LatLng{}, fmt.Errorf("could not parse %s: %w", match[2], err)
	}

	return LatLng{Lat: lat, Lng: lng}, nil
}
