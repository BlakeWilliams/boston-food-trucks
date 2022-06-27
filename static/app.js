const today = new Date().toLocaleDateString(undefined, { weekday: "long" });
const allMarkers = [];
let map;

document.addEventListener(
  "DOMContentLoaded",
  function () {
    map = new mapboxgl.Map({
      container: "map",
      style: "mapbox://styles/mapbox/streets-v11",
      center: [-71.0675, 42.3555],
      zoom: 14,
      style: window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "mapbox://styles/mapbox/dark-v10"
        : "mapbox://styles/mapbox/streets-v11",
    });

    for (const key of Object.keys(truckData)) {
      const trucks = truckData[key];
      const latlng = [trucks[0].Lng, trucks[0].Lat];

      const marker = new mapboxgl.Marker({ draggable: false })
        .setLngLat(latlng)
        .addTo(map);
      allMarkers.push(marker);

      const link = document.querySelector(`[data-location="${key}"]`);
      link.addEventListener("click", () => {
        onLocationClick(marker, latlng);
      });

      const popup = createPopup(trucks);
      marker.setPopup(popup);
    }
  },
  false
);

function createPopup(trucks) {
  const trucksAvailableToday = trucks.filter(
    (truck) => !!truck.Schedule[today]
  );
  const trucksNotAvailableToday = trucks.filter(
    (truck) => !truck.Schedule[today]
  );

  let popupHTML = `<div class="text-gray-800">`;

  if (trucksAvailableToday.length > 0) {
    popupHTML += textForTrucks("Unavailable", trucksAvailableToday);
  }

  if (trucksNotAvailableToday.length > 0) {
    popupHTML += textForTrucks("Unavailable", trucksNotAvailableToday);
  }

  popupHTML += `</div>`;
  return new mapboxgl.Popup().setHTML(popupHTML);
}

function textForTrucks(label, trucks) {
  let output = "";
  if (output != "") {
    output += `<div class="mt-3"></div>`;
  }

  output += `<h1 class="text-lg mb-1">${label}</h1><ul class="list-disc pl-4">`;
  trucks.forEach((truck) => {
    output += `<li class="text-base mb-1">`;
    if (truck.Schedule[today]) {
      output += `${truck.Name} - ${truck.Schedule[today]}`;
    } else {
      output += truck.Name;
    }
    output += `</li>`;
  });
  output += `</ul>`;

  return output;
}

function onLocationClick(marker, latlng) {
  map.panTo(latlng);

  // close all popups
  allMarkers.forEach((marker) => {
    if (marker.getPopup().isOpen()) {
      marker.togglePopup();
    }
  });

  if (!marker.getPopup().isOpen()) {
    marker.togglePopup();
  }
}
