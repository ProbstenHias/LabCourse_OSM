var map = L.map('map', {
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
let startIdx
let end
let endIdx
let line

function setPoint(e, isStart) {
    let point = L.marker()
    point.setLatLng(e.latlng)
    if (isStart) {
        fetchPoint(point, isStart)
    } else {
        fetchPoint(point, isStart).then(fetchRoute);
    }

}


function onEachLineFeature(feature, layer) {
    if (feature.properties && feature.properties.popupContent) {
        layer.bindPopup(feature.properties.popupContent);
    }
}

function onEachStartFeature(feature, layer) {
    layer.setIcon(greenIcon)
    if (feature.properties && feature.properties.popupContent) {
        layer.bindPopup(feature.properties.popupContent);
    }
    if (feature.properties && feature.properties.index) {
        startIdx = feature.properties.index
    }
}

function onEachEndFeature(feature, layer) {
    layer.setIcon(redIcon)
    if (feature.properties && feature.properties.popupContent) {
        layer.bindPopup(feature.properties.popupContent);
    }
    if (feature.properties && feature.properties.index) {
        endIdx = feature.properties.index
    }
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
        setPoint(e, true);
    } else if (end == null) {
        setPoint(e, false);
    } else {
        map.removeLayer(end);
        end = null;
        map.removeLayer(start)
        start = null
        map.removeLayer(line)
        line = null
        setPoint(e, true);
    }
}

map.on("click", OnMapClick)


function fetchPoint(point, isStart) {
    const data = {
        lat: point.getLatLng().lat, lng: point.getLatLng().lng
    };
    //build URL
    let url = new URL("http://localhost:8081/point");
    for (let k in data) {
        url.searchParams.append(k, data[k]);
    }
    console.log("fetching data from url: " + url)

    return fetch(url)

        .then((result) => {
            console.log(result);
            if (result.status !== 200) {
                throw new Error("Bad Server Response");
            }
            return result.text();
        })

        .then((response) => {
            let fc = JSON.parse(response)
            if (isStart) {
                start = L.geoJSON(fc, {
                    onEachFeature: onEachStartFeature
                })
                // start.setIcon(greenIcon)
                start.addTo(map)

            } else {
                end = L.geoJSON(fc, {
                    onEachFeature: onEachEndFeature
                })
                // end.setIcon(redIcon)
                end.addTo(map)
            }
        })

        .catch((error) => {
            console.log(error);
        });
}

function fetchRoute() {
    const data = {
        startIdx: startIdx, endIdx: endIdx
    };

    // build URL
    let url = new URL("http://localhost:8081/route");
    for (let k in data) {
        url.searchParams.append(k, data[k]);
    }
    console.log("fetching data from url: " + url)

    fetch(url)

        .then((result) => {
            console.log(result);
            if (result.status !== 200 && result.status !== 204) {
                throw new Error("Bad Server Response");
            }
            if (result.status === 204) {
                throw new Error("No Path Found")
            }
            return result.text();
        })

        .then((response) => {
            line = L.geoJSON(JSON.parse(response), {
                onEachFeature: onEachLineFeature
            });
            line.addTo(map)
        })

        .catch((error) => {
            console.log(error);
        });
}
