var map = L.map('map',{
    minZoom: 3
})
map.setView([0, 0], 0);
L.tileLayer('https://api.maptiler.com/maps/streets/{z}/{x}/{y}.png?key=AWTkDE2AmJ3tv0oa3YoE', {
    attribution: "<a href=\"https://www.maptiler.com/copyright/\" target=\"_blank\">&copy; MapTiler</a> <a href=\"https://www.openstreetmap.org/copyright\" target=\"_blank\">&copy; OpenStreetMap contributors</a>"
}).addTo(map);


var greenIcon = new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

var redIcon = new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});


let start
let end

let line

function setStartingPoint(e) {
    start = L.marker();
    start
        .setLatLng(e.latlng)
        .setIcon(greenIcon)
        .addTo(map)
        .on("click", onHoverMarker(start))
}

function setEndingPoint(e) {
    end = L.marker();
    end
        .setLatLng(e.latlng)
        .setIcon(redIcon)
        .addTo(map)
        .on("click", onHoverMarker(end))
}


function onHoverMarker(point) {
    return point.bindPopup(`Lat: ${point.getLatLng().lat}<br>Lng: ${point.getLatLng().lng}`).openPopup();
}

function OnMapClick(e) {
    if (e.latlng.lng < -180) {
        alert("That's too far to the left!");
        return;
    }
    if (e.latlng.lng > 180) {
        alert("That's too far to the right!");
        return;
    }
    if (start == null) {
        setStartingPoint(e);
    } else if (end == null) {
        setEndingPoint(e);
        fetchRoute();
    } else {
        map.removeLayer(end);
        end = null;
        map.removeLayer(start)
        start = null
        map.removeLayer(line)
        line = null
        setStartingPoint(e);
    }
}

map.on("click", OnMapClick)

function fetchRoute() {
    const data = {
        startLat: start.getLatLng().lat,
        startLng: start.getLatLng().lng,
        endLat: end.getLatLng().lat,
        endLng: end.getLatLng().lng
    };

    // build URL
    let url = new URL("http://localhost:8080/route");
    for (let k in data) {
        url.searchParams.append(k, data[k]);
    }
    console.log("fetching data from url: " + url)

    fetch(url)

        .then((result) => {
            console.log(result);
            if (result.status !== 200) {
                throw new Error("Bad Server Response");
            }
            return result.text();
        })

        .then((response) => {
            console.log(response);
            line = L.geoJSON(JSON.parse(response));
            line.addTo(map)
        })

        .catch((error) => {
            console.log(error);
        });
}
